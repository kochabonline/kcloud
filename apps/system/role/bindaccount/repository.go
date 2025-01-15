package bindaccount

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Bind 绑定账户角色
func (repo *Repository) Create(ctx context.Context, roleBindAccount *RoleBindAccount) error {
	return repo.db.WithContext(ctx).Create(roleBindAccount).Error
}

// Unbind 解绑账户角色
func (repo *Repository) Delete(ctx context.Context, roleBindAccount *RoleBindAccount) error {
	return repo.db.WithContext(ctx).Where("account_id = ? AND role_id = ?", roleBindAccount.AccountId, roleBindAccount.RoleId).Delete(roleBindAccount).Error
}

// 根据账户id查询角色
func (repo *Repository) FindByAccountId(ctx context.Context, accountId int64) ([]*RoleBindAccount, error) {
	var roleBindAccounts []*RoleBindAccount
	if err := repo.db.WithContext(ctx).Where("account_id = ?", accountId).Find(&roleBindAccounts).Error; err != nil {
		return nil, err
	}
	return roleBindAccounts, nil
}
