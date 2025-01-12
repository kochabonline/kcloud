package audit

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, audit *Audit) error {
	return r.db.WithContext(ctx).Create(audit).Error
}

func (r *Repository) FindByAccountId(ctx context.Context, accountId int64) (*Audit, error) {
	var audit Audit
	err := r.db.WithContext(ctx).Where("account_id = ?", accountId).First(&audit).Error
	return &audit, err
}

func (r *Repository) Update(ctx context.Context, audit *Audit) error {
	return r.db.WithContext(ctx).Save(audit).Error
}
