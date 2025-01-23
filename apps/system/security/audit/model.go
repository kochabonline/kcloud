package audit

import "github.com/kochabonline/kcloud/apps/common"

type Audit struct {
	common.Meta
	// 账户id
	AccountId int64 `json:"accountId" gorm:"index;not null;type:bigint;comment:账户id"`
	// 账户名称
	AccountUsername string `json:"accountUsername" gorm:"index;not null;type:varchar(32);comment:账户名称"`
	// 最近登录时间
	LastLoginTimestamp int64 `json:"lastLoginTimestamp" gorm:"not null;type:bigint;comment:最近登录时间"`
	// 最近登录ip
	LastLoginIp string `json:"lastLoginIp" goem:"not null;type:varchar(32);comment:最近登录ip"`
	// 最近登录地点
	LastLoginLocation string `json:"lastLoginLocation" gorm:"not null;type:varchar(64);comment:最近登录地点"`
	// 最近登录设备
	LastLoginDevice string `json:"lastLoginDevice" gorm:"not null;type:varchar(64);comment:最近登录设备"`
	// 最近登录ua
	LastLoginUserAgent string `json:"lastLoginUserAgent" gorm:"not null;type:varchar(256);comment:最近登录ua"`
}

func (m *Audit) TableName() string {
	return "audit"
}
