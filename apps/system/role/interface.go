package role

import (
	"context"
)

type Interface interface {
	// 创建角色
	Create(ctx context.Context, req *Request) error
	// 更新角色
	Update(ctx context.Context, req *Request) error
	// 删除角色
	Delete(ctx context.Context, req *DeleteRequest) error
	// 角色列表
	FindAll(ctx context.Context, req *FindAllRequest) (*Roles, error)
	// 根据角色列表获取角色列表路由
	FindRouteByRoles(ctx context.Context, req *FindRouteByRolesRequest) (*FindRouteByRolesResponse, error)
}

type Request struct {
	// 角色名
	Name string `json:"name" validate:"required,max=24"`
	// 角色代码
	Code string `json:"code" validate:"required,max=24"`
	// 状态 1:启用 2:禁用
	Status int `json:"status" validate:"required,oneof=1 2"`
	// 描述
	Description string `json:"description" validate:"max=128"`
}

type DeleteRequest struct {
	// 角色Id
	Id int64 `uri:"id" validate:"required"`
}

type FindAllRequest struct {
	// 页码
	Page int `form:"page" validate:"omitempty,min=1"`
	// 每页数量
	Size int `form:"size" validate:"omitempty,min=1,max=100"`
}

type Roles struct {
	// 总数
	Total int64 `json:"total"`
	// 角色列表
	Items []*Role `json:"items"`
}

type FindRouteByRolesRequest struct {
	// 角色列表
	Roles []string `json:"roles" validate:"required"`
}

type FindRouteByRolesResponse struct {
	// 路由列表
	Routes []string `json:"routes"`
}
