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
	tables := []any{account.Account{}, google.GoogleAuth{}, audit.Audit{}, device.Device{}}
	if err := m.db.AutoMigrate(tables...); err != nil {
		log.Fatalw("message", "自动迁移数据库失败", "error", err.Error())
	}
	log.Info("自动迁移数据库成功")
}

// init 初始化数据
func (m *Migrate) init() {
	m.initMongoIndex()
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

// initadmin 初始化管理员账户
func (m *Migrate) initadmin() {
	// 初始化管理员账户
	now := time.Now().Unix()
	password := uuid.New().String()
	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		log.Fatalw("message", "密码加密失败", "error", err.Error())
	}
	admin := account.Account{
		Meta: common.Meta{
			CreatedAt: now,
		},
		Username: "admin",
		Password: hashedPassword,
		Role:     common.RoleAdmin,
		Status:   common.StatusNormal,
	}
	// 查询是否存在
	if err := m.db.Where("username = ?", admin.Username).First(&account.Account{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		if err := m.db.Create(&admin).Error; err != nil {
			log.Fatalw("message", "初始化管理员账户失败", "error", err.Error())
		}
		log.Infow("message", "初始化管理员账户成功", "username", admin.Username, "password", password)
	}
}
