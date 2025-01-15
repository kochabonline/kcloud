package role

import (
	"context"

	"github.com/kochabonline/kcloud/internal/util"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, role *Role) error {
	return repo.db.WithContext(ctx).Create(role).Error
}

func (repo *Repository) Update(ctx context.Context, role *Role) error {
	return repo.db.WithContext(ctx).Model(role).Updates(role).Error
}

func (repo *Repository) Delete(ctx context.Context, id int64) error {
	return repo.db.WithContext(ctx).Delete(&Role{}, id).Error
}

func (repo *Repository) FindAll(ctx context.Context, req *FindAllRequest) (*Roles, error) {
	var roles Roles
	query := repo.db.WithContext(ctx)

	offset, limit := util.Paginate(req.Page, req.Size)
	if err := query.Preload("Menus").Count(&roles.Total).Offset(offset).Limit(limit).Find(&roles.Items).Error; err != nil {
		return nil, err
	}

	return &roles, nil
}
