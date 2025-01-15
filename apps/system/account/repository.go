package account

import (
	"context"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/internal/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db       *gorm.DB
	rediscli *redis.Client
}

func NewRepository(db *gorm.DB, rediscli *redis.Client) *Repository {
	return &Repository{
		db:       db,
		rediscli: rediscli,
	}
}

func (repo *Repository) Create(ctx context.Context, account Account) error {
	return repo.db.WithContext(ctx).Create(&account).Error
}

func (repo *Repository) ChangePassword(ctx context.Context, id int64, password string) error {
	return repo.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("password", password).Error
}

func (repo *Repository) ChangeEmail(ctx context.Context, id int64, email string) error {
	return repo.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("email", email).Error
}

func (repo *Repository) ChangeMobile(ctx context.Context, id int64, mobile string) error {
	return repo.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("mobile", mobile).Error
}

func (repo *Repository) FindByUsername(ctx context.Context, username string) (*Account, error) {
	var account Account
	err := repo.db.WithContext(ctx).
		Preload("Roles").
		Where("username = ?", username).
		Where("status = ?", common.StatusNormal).
		First(&account).Error
	return &account, err
}

func (repo *Repository) FindById(ctx context.Context, id int64) (*Account, error) {
	var account Account
	err := repo.db.WithContext(ctx).
		Preload("Roles").
		Where("id = ?", id).
		Where("status = ?", common.StatusNormal).
		First(&account).Error
	return &account, err
}

func (repo *Repository) FindAll(ctx context.Context, req *FindAllRequest) (*Accounts, error) {
	var accounts Accounts

	// 根据请求参数构建查询条件
	query := repo.db.WithContext(ctx).
		Model(&Account{}).
		Preload("Roles")

	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = ?", common.StatusNormal)
	}
	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}
	if req.Keyword != "" {
		query = query.Where("username like ?", "%"+req.Keyword+"%")
	}

	offset, limit := util.Paginate(req.Page, req.Size)
	if err := query.Count(&accounts.Total).Offset(offset).Limit(limit).Find(&accounts.Items).Error; err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (repo *Repository) Delete(ctx context.Context, id int64) error {
	// 开始事务
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除角色关联
	if err := tx.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 软删除
	if err := tx.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("status", common.StatusDeleted).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
