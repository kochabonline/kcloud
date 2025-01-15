package bindaccount

import (
	"context"

	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo *Repository
	log  log.Helper
}

func NewController(repo *Repository, log log.Helper) *Controller {
	return &Controller{
		repo: repo,
		log:  log,
	}
}

// Bind 绑定账户角色
func (ctrl *Controller) Bind(ctx context.Context, req *request) error {
	roleBindAccount := &RoleBindAccount{
		AccountId: req.AccountId,
		RoleId:    req.RoleId,
	}

	if err := ctrl.repo.Create(ctx, roleBindAccount); err != nil {
		ctrl.log.Errorw("message", "绑定账户角色失败", "request", req, "error", err.Error())
		return err
	}

	return nil
}

// Unbind 解绑账户角色
func (ctrl *Controller) Unbind(ctx context.Context, req *request) error {
	roleBindAccount := &RoleBindAccount{
		AccountId: req.AccountId,
		RoleId:    req.RoleId,
	}

	if err := ctrl.repo.Delete(ctx, roleBindAccount); err != nil {
		ctrl.log.Errorw("message", "解绑账户角色失败", "request", req, "error", err.Error())
		return err
	}

	return nil
}

// 根据账户id查询角色
func (ctrl *Controller) FindByAccountId(ctx context.Context, req *findByAccountIdRequest) ([]*RoleBindAccount, error) {
	var roleBindAccounts []*RoleBindAccount
	roleBindAccounts, err := ctrl.repo.FindByAccountId(ctx, req.AccountId)
	if err != nil {
		ctrl.log.Errorw("message", "根据账户id查询角色失败", "request", req, "error", err.Error())
		return roleBindAccounts, err
	}

	return roleBindAccounts, nil
}
