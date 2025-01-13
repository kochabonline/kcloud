package menu

import (
	"context"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/core/convert/tree"
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

func (ctrl *Controller) Create(ctx context.Context, req *Request) error {
	menu := &Menu{
		Meta: common.Meta{
			CreatedAt: time.Now().Unix(),
		},
		Name:        req.Name,
		ParentId:    req.ParentId,
		Icon:        req.Icon,
		Router:      req.Router,
		Component:   req.Component,
		Title:       req.Title,
		I18n:        req.I18n,
		Status:      req.Status,
		HideMenu:    req.HideMenu,
		KeeperAlive: req.KeeperAlive,
		Order:       req.Order,
	}

	if err := ctrl.repo.Create(ctx, menu); err != nil {
		ctrl.log.Errorw("message", "新增菜单失败", "error", err.Error())
		return err
	}

	return nil
}

// 修改菜单
func (ctrl *Controller) Update(ctx context.Context, req *Request) error {
	menu := &Menu{
		Meta: common.Meta{
			CreatedAt: time.Now().Unix(),
		},
		Name:        req.Name,
		ParentId:    req.ParentId,
		Icon:        req.Icon,
		Router:      req.Router,
		Component:   req.Component,
		Title:       req.Title,
		I18n:        req.I18n,
		Status:      req.Status,
		HideMenu:    req.HideMenu,
		KeeperAlive: req.KeeperAlive,
		Order:       req.Order,
	}

	if err := ctrl.repo.Update(ctx, menu); err != nil {
		ctrl.log.Errorw("message", "修改菜单失败", "error", err.Error())
		return err
	}

	return nil
}

// 删除菜单
func (ctrl *Controller) Delete(ctx context.Context, req *DeleteRequest) error {
	if err := ctrl.repo.Delete(ctx, req.Id); err != nil {
		ctrl.log.Errorw("message", "删除菜单失败", "error", err.Error())
		return err
	}
	return nil
}

// 树形菜单
func (ctrl *Controller) Tree(ctx context.Context) ([]tree.Node, error) {
	menus, err := ctrl.repo.FindAll(ctx)
	if err != nil {
		ctrl.log.Errorw("message", "查询菜单列表失败", "error", err.Error())
		return nil, err
	}
	return tree.BuildTree(tree.ConvertToNodeSlice(menus), 0), nil
}
