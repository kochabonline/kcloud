package bindaccount

import "context"

type Interface interface {
	// 绑定账户角色
	Bind(ctx context.Context, req *request) error
	// 解绑账户角色
	Unbind(ctx context.Context, req *request) error
	// 根据账户id查询角色
	FindByAccountId(ctx context.Context, req *findByAccountIdRequest) ([]*RoleBindAccount, error)
}

type request struct {
	// 账户id
	AccountId int64 `json:"accountId" validate:"required" label:"账户id"`
	// 角色id
	RoleId int64 `json:"roleId" validate:"required" label:"角色id"`
}

type findByAccountIdRequest struct {
	// 账户id
	AccountId int64 `uri:"account_id" validate:"required" label:"账户id"`
}
