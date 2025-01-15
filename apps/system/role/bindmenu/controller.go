package bindmenu

import (
	"context"

	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo *Repository
	log log.Helper
}

func NewController(repo *Repository, log log.Helper) *Controller {
	return &Controller{
		repo: repo,
		log:  log,
	}
}

func (ctrl *Controller) Edit(ctx context.Context, req *Request) error {
	if err := ctrl.repo.Edit(ctx, req.RoleId, req.MenuIds); err != nil {
		ctrl.log.Errorw("message", "编辑角色菜单失败", "error", err.Error())
		return err
	}

	return nil
}