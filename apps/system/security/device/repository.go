package device

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

func (repo *Repository) Create(ctx context.Context, device *Device) error {
	return repo.db.WithContext(ctx).Create(device).Error
}

func (repo *Repository) Update(ctx context.Context, device *Device) error {
	return repo.db.WithContext(ctx).Model(device).Updates(device).Error
}

func (repo *Repository) Delete(ctx context.Context, device *Device) error {
	return repo.db.WithContext(ctx).Where("account_id AND device_id = ?", device.AccountId, device.DeviceId).Delete(&Device{}).Error
}

func (repo *Repository) FindAll(ctx context.Context, req *FindAllRequest) (*Devices, error) {
	var devices Devices
	query := repo.db.WithContext(ctx)
	if req.AccountId != "" {
		query = query.Where("account_id = ?", req.AccountId)
	}
	if req.DeviceId != "" {
		query = query.Where("device_id = ?", req.DeviceId)
	}
	if err := query.Model(&Device{}).Count(&devices.Total).Error; err != nil {
		return nil, err
	}
	offset, limit := util.Paginate(req.Page, req.Size)
	if err := query.Offset(offset).Limit(limit).Find(&devices.Items).Error; err != nil {
		return nil, err
	}

	return &devices, nil
}

func (repo *Repository) FindByAccountId(ctx context.Context, accountId string) (*Device, error) {
	var device Device
	if err := repo.db.WithContext(ctx).Where("account_id = ?", accountId).First(&device).Error; err != nil {
		return nil, err
	}

	return &device, nil
}
