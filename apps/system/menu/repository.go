package menu

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

func (repo *Repository) Create(ctx context.Context, menu *Menu) error {
	return repo.db.WithContext(ctx).Create(menu).Error
}

func (repo *Repository) Update(ctx context.Context, menu *Menu) error {
	return repo.db.WithContext(ctx).Model(menu).Updates(menu).Error
}

func (repo *Repository) Delete(ctx context.Context, id int64) error {
	return repo.db.WithContext(ctx).Delete(&Menu{}, id).Error
}

func (repo *Repository) FindAll(ctx context.Context) ([]*Menu, error) {
	var menus []*Menu
	err := repo.db.WithContext(ctx).Find(&menus).Error
	return menus, err
}