package device

import "github.com/kochabonline/kcloud/apps/common"

type Device struct {
	common.Meta
	// 账户id
	AccountId string `json:"account_id" gorm:"index;not null;type:bigint;comment:账户id"`
	// 账户名称
	AccountName string `json:"account_name" gorm:"index;not null;type:varchar(32);comment:账户名称"`
	// 设备id
	DeviceId string `json:"device_id" gorm:"index;not null;type:varchar(32);comment:设备id"`
	// 设备名称
	Name string `json:"name" gorm:"not null;type:varchar(32);comment:设备名称"`
	// 设备类型
	Type string `json:"type" gorm:"not null;type:varchar(32);comment:设备类型"`
}

func (m *Device) TableName() string {
	return "device"
}
