package account

import (
	"context"

	"github.com/kochabonline/kcloud/apps/common"
)

type Interface interface {
	// 创建用户
	Create(ctx context.Context, req *CreateRequest) error
	// 更改密码
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) error
	// 更改邮箱
	ChangeEmail(ctx context.Context, req *ChangeEmailRequest) error
	// 更改手机号
	ChangeMobile(ctx context.Context, req *ChangeMobileRequest) error
	// 查找所有用户
	FindAll(ctx context.Context, req *FindAllRequest) (*Accounts, error)
	// 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*Account, error)
	// 账户详情(手机号、邮箱等脱敏返回)
	Detail(ctx context.Context) (*Account, error)
	// 删除账户
	Delete(ctx context.Context, req *DeleteRequest) error
}

// 创建用户请求
type CreateRequest struct {
	// 用户名
	Username string `json:"username" validate:"required,min=3,max=24" label:"用户名"`
	// 密码
	Password string `json:"password" validate:"required" label:"密码"`
}

// 修改密码请求
type ChangePasswordRequest struct {
	// 新密码
	Password string `json:"password" validate:"required" label:"新密码"`
	// 安全码
	SecurityCode string `json:"securityCode" validate:"required" label:"安全码"`
}

// 修改邮箱请求
type ChangeEmailRequest struct {
	// 新邮箱
	Email string `json:"email" validate:"required,email" label:"邮箱"`
	// 安全码
	SecurityCode string `json:"securityCode" validate:"required" label:"安全码"`
}

// 修改手机号请求
type ChangeMobileRequest struct {
	// 新手机号
	Mobile string `json:"mobile" validate:"required" label:"手机号"`
	// 安全码
	SecurityCode string `json:"securityCode" validate:"required" label:"安全码"`
}

// 用户列表查询请求
type FindAllRequest struct {
	// 页码
	Page int `form:"page" validate:"omitempty,gt=0" label:"页码"`
	// 每页数量
	Size int `form:"size" validate:"omitempty,gt=0" label:"每页数量"`
	// 状态
	Status common.Status `form:"status" validate:"omitempty" label:"状态"`
	// 用户名
	Username string `form:"username" validate:"omitempty" label:"用户名"`
	// 关键字
	Keyword string `form:"keyword" validate:"omitempty" label:"关键字"`
}

// 用户列表查询响应
type Accounts struct {
	// 总数
	Total int64 `json:"total"`
	// 用户列表
	Items []*Account `json:"items"`
}

// 删除用户请求
type DeleteRequest struct {
	// 用户id
	Id int64 `uri:"id" validate:"required" label:"用户id"`
}
