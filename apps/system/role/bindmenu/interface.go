package bindmenu

import "context"

type Interface interface {
	// 编辑角色菜单, 可绑定解绑多个菜单
	Edit(ctx context.Context, req *Request) error
}

type Request struct {
	// 角色Id
	RoleId int64 `json:"roleId" validate:"required" comment:"角色Id"`
	// 菜单Id
	MenuIds []int64 `json:"menuIds" validate:"required" comment:"菜单Id"`
}
