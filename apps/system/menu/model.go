package menu

import (
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/core/util/tree"
)

type Menu struct {
	common.Meta
	// 父级菜单id
	ParentId int64 `json:"parent_id" gorm:"type:tinyint(1);default:0;index;comment:父级菜单id"`
	// 菜单类型
	Type int `json:"type" gorm:"type:tinyint(1);default:0;comment:菜单类型"`
	// 菜单名称
	Name string `json:"name" gorm:"type:varchar(50);comment:菜单名称"`
	// 路由名称
	Router string `json:"router" gorm:"type:varchar(50);comment:路由名称"`
	// 路由路径
	Component string `json:"component" gorm:"type:varchar(50);comment:路由路径"`
	// 菜单状态
	Status int `json:"status" gorm:"type:tinyint(1);comment:菜单状态"`
	// 菜单标题
	Title string `json:"title" gorm:"type:varchar(50);comment:标题"`
	// 国际化
	I18n string `json:"i18n" gorm:"type:varchar(50);comment:国际化"`
	// 隐藏菜单
	HideMenu bool `json:"hide_menu" gorm:"comment:隐藏菜单"`
	// 缓存菜单
	KeeperAlive bool `json:"keeper_alive" gorm:"comment:缓存菜单"`
	// 菜单图标
	Icon string `json:"icon" gorm:"type:varchar(50);comment:图标"`
	// 菜单排序
	Order int `json:"order" gorm:"type:int(11);comment:排序"`
	// 子菜单
	Children []*Menu `json:"children" gorm:"-"`
}

func (Menu) TableName() string {
	return "menu"
}

func (m *Menu) GetNode() tree.Node {
	return m
}
func (m *Menu) SetChildren(children []tree.Node) {
	m.Children = make([]*Menu, len(children))
	for i, child := range children {
		m.Children[i] = child.(*Menu)
	}
}
func (m *Menu) GetId() int64 {
	return m.Id
}
func (m *Menu) GetParentId() int64 {
	return m.ParentId
}
