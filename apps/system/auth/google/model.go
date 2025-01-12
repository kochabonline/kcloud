package google

import "github.com/kochabonline/kcloud/apps/common"

type GoogleAuth struct {
	common.Meta
	// 账户Id
	AccountId int64 `json:"account_id" gorm:"not null;index;comment:账户Id"`
	// 谷歌验证器密钥
	Secret string `json:"secret" gorm:"not null;comment:谷歌验证器密钥"`
	// 状态
	Status common.Status `json:"status" gorm:"type:tinyint(1);default:2;comment:状态 1:正常 2:禁用 3:删除"`
}

func (GoogleAuth) TableName() string {
	return "auth_google"
}
