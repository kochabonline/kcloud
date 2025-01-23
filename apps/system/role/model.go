package role

import (
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/menu"
)

type Role struct {
	common.Meta
	// 角色名
	Name string `json:"name" gorm:"type:varchar(24);not null;uniqueIndex:idx_deleted_at;comment:角色名"`
	// 角色编码
	Code string `json:"code" gorm:"type:varchar(24);not null;uniqueIndex:idx_deleted_at;comment:角色编码"`
	// 状态
	Status int `json:"status" gorm:"type:tinyint(1);default:1;comment:状态 1:启用 2:禁用"`
	// 描述
	Description string `json:"description" gorm:"type:varchar(128);comment:描述"`
	// 菜单
	Menus []menu.Menu `json:"-" gorm:"many2many:role_bind_menu;"`
}

func (Role) TableName() string {
	return "role"
}
