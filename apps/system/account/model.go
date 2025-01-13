package account

import (
	"github.com/kochabonline/kcloud/apps/common"
)

type Account struct {
	common.Meta
	// 用户名
	// uniqueIndex:idx_deleted_at: 唯一索引, 避免软删除的数据重复
	Username string `json:"username" gorm:"type:varchar(24);not null;uniqueIndex:idx_deleted_at;comment:用户名"`
	// 密码
	Password string `json:"-" gorm:"type:varchar(64);not null;comment:密码"`
	// 昵称
	Nickname string `json:"nickname" gorm:"type:varchar(24);comment:昵称"`
	// 邮箱
	Email string `json:"email" gorm:"type:varchar(48);index;comment:邮箱"`
	// 手机号
	Mobile string `json:"mobile" gorm:"type:varchar(24);index;comment:手机号"`
	// 角色
	Role common.Role `json:"role" gorm:"type:tinyint(1);default:1;comment:角色 1:普通用户 2:管理员"`
	// 状态
	Status common.Status `json:"status" gorm:"type:tinyint(1);default:1;comment:状态 1:正常 2:禁用 3:删除"`
}

func (Account) TableName() string {
	return "account"
}
