package migrate

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kcloud/apps/system/auth/google"
	"github.com/kochabonline/kcloud/apps/system/menu"
	"github.com/kochabonline/kcloud/apps/system/role"
	"github.com/kochabonline/kcloud/apps/system/role/bindaccount"
	"github.com/kochabonline/kcloud/apps/system/role/bindmenu"
	"github.com/kochabonline/kcloud/apps/system/security/audit"
	"github.com/kochabonline/kcloud/apps/system/security/device"
	"github.com/kochabonline/kit/core/crypto/bcrypt"
	"github.com/kochabonline/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"gorm.io/gorm"
)

type Migrate struct {
	db       *gorm.DB
	mongocli *mongo.Client
}

func NewMigrate(db *gorm.DB, client *mongo.Client) *Migrate {
	return &Migrate{
		db:       db,
		mongocli: client,
	}
}

func (m *Migrate) Start() {
	m.migrate()
	m.init()
}

// migrate 数据库迁移
func (m *Migrate) migrate() {
	tables := []any{
		account.Account{},
		role.Role{},
		menu.Menu{},
		bindaccount.RoleBindAccount{},
		bindmenu.RoleBindMenu{},
		google.GoogleAuth{},
		audit.Audit{},
		device.Device{},
	}
	if err := m.db.AutoMigrate(tables...); err != nil {
		log.Fatalw("message", "自动迁移数据库失败", "error", err.Error())
	}
	log.Info("自动迁移数据库成功")
}

// init 初始化数据
func (m *Migrate) init() {
	m.initMongoIndex()
	m.initRole()
	m.initadmin()
}

// initMongoIndex 初始化Mongo索引
func (m *Migrate) initMongoIndex() {
	channalColl := m.mongocli.Database("kcloud").Collection("channal")
	models := []mongo.IndexModel{
		{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"api_key": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"created_at": -1}},
		{Keys: bson.M{"updated_at": -1}},
	}

	names, err := channalColl.Indexes().CreateMany(context.TODO(), models)
	if err != nil {
		log.Fatalw("message", "创建索引失败", "error", err.Error())
	}
	log.Infow("message", "创建索引成功", "indexs", strings.Join(names, ","))

	messageColl := m.mongocli.Database("kcloud").Collection("message")
	models = []mongo.IndexModel{
		{Keys: bson.M{"channal_id": 1}},
		{Keys: bson.M{"created_at": -1}},
		{Keys: bson.M{"updated_at": -1}},
	}

	names, err = messageColl.Indexes().CreateMany(context.TODO(), models)
	if err != nil {
		log.Fatalw("message", "创建索引失败", "error", err.Error())
	}
	log.Infow("message", "创建索引成功", "indexs", strings.Join(names, ","))
}

// initRole 初始化角色
func (m *Migrate) initRole() {
	now := time.Now().Unix()
	roles := []role.Role{
		{
			Meta: common.Meta{
				CreatedAt: now,
			},
			Name:        common.RoleNormal,
			Code:        "R_NORMAL",
			Description: "普通账户",
			Status:      int(common.StatusNormal),
		},
		{
			Meta: common.Meta{
				CreatedAt: now,
			},
			Name:        common.RoleAdmin,
			Code:        "R_ADMIN",
			Description: "管理员账户",
			Status:      int(common.StatusNormal),
		},
		{
			Meta: common.Meta{
				CreatedAt: now,
			},
			Name:        common.RoleSuperAdmin,
			Code:        "R_SUPER",
			Description: "超级管理员账户",
			Status:      int(common.StatusNormal),
		},
	}
	for _, r := range roles {
		if err := m.db.Where("name = ?", r.Name).First(&role.Role{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			if err := m.db.Create(&r).Error; err != nil {
				log.Fatalw("message", "初始化角色失败", "role", r, "error", err.Error())
			}
		}
	}
	log.Infow("message", "初始化角色成功", "roles", strings.Join([]string{common.RoleNormal, common.RoleAdmin, common.RoleSuperAdmin}, ","))
}

// initadmin 初始化超级管理员账户
func (m *Migrate) initadmin() {
	// 初始化超级管理员账户
	now := time.Now().Unix()
	password := uuid.New().String()
	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		log.Fatalw("message", "密码加密失败", "error", err.Error())
	}
	superadmin := account.Account{
		Meta: common.Meta{
			CreatedAt: now,
		},
		Username: "superadmin",
		Password: hashedPassword,
		Status:   common.StatusNormal,
	}

	// 查询超级管理员角色和是否存在管理员账户
	var role role.Role
	var account account.Account
	if err := m.db.Where("name = ?", common.RoleSuperAdmin).First(&role).Error; err != nil {
		log.Fatalw("message", "查询超级管理员角色失败", "error", err.Error())
	}
	if err := m.db.Where("username = ?", superadmin.Username).First(&account).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalw("message", "查询超级管理员角色失败", "error", err.Error())
	}
	if account.Id != 0 {
		return
	}

	// 关联超级管理员角色
	superadmin.Roles = append(superadmin.Roles, role)

	// 创建用户并关联角色
	if err := m.db.Create(&superadmin).Error; err != nil {
		log.Fatalw("message", "初始化管理员账户失败", "error", err.Error())
	}

	log.Infow("message", "初始化管理员账户成功", "username", superadmin.Username, "password", password)
}
