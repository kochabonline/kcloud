package menu

import (
	"context"

	"github.com/kochabonline/kit/core/util/tree"
)

type Interface interface {
	// 新增菜单
	Create(ctx context.Context, req *Request) error
	// 修改菜单
	Update(ctx context.Context, req *Request) error
	// 删除菜单
	Delete(ctx context.Context, req *DeleteRequest) error
	// 树形菜单
	Tree(ctx context.Context) ([]tree.Node, error)
}

type Request struct {
	// 父级菜单id
	ParentId int64 `json:"parent_id" validate:"required" label:"父级菜单id"`
	// 菜单类型
	Type int `json:"type" validate:"required" label:"菜单类型"`
	// 菜单名称
	Name string `json:"name" validate:"required" label:"菜单名称"`
	// 路由名称
	Router string `json:"router" validate:"required" label:"路由名称"`
	// 路由路径
	Component string `json:"component" validate:"required" label:"路由路径"`
	// 菜单状态
	Status int `json:"status" validate:"required" label:"菜单状态"`
	// 菜单标题
	Title string `json:"title" validate:"required" label:"标题"`
	// 国际化
	I18n string `json:"i18n" validate:"omitempty" label:"国际化"`
	// 隐藏菜单
	HideMenu bool `json:"hide_menu" validate:"omitempty" label:"隐藏菜单"`
	// 缓存菜单
	KeeperAlive bool `json:"keeper_alive" validate:"omitempty" label:"缓存菜单"`
	// 菜单图标
	Icon string `json:"icon" validate:"required" label:"图标"`
	// 菜单排序
	Order int `json:"order" validate:"omitempty" label:"排序"`
}

type DeleteRequest struct {
	Id int64 `uri:"id" validate:"required" label:"菜单id"`
}
