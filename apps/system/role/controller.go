package role

import (
	"context"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/core/util/slice"
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
	// 角色名称不能为系统角色
	if req.Name == common.RoleNormal ||
		req.Name == common.RoleAdmin ||
		req.Name == common.RoleSuperAdmin {
		return common.ErrorInvalidRole
	}

	role := &Role{
		Name:        req.Name,
		Status:      req.Status,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := ctrl.repo.Create(ctx, role); err != nil {
		ctrl.log.Errorw("message", "新增角色失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) Update(ctx context.Context, req *Request) error {
	role := &Role{
		Name:        req.Name,
		Status:      req.Status,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := ctrl.repo.Update(ctx, role); err != nil {
		ctrl.log.Errorw("message", "更新角色失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) Delete(ctx context.Context, req *DeleteRequest) error {
	if err := ctrl.repo.Delete(ctx, req.Id); err != nil {
		ctrl.log.Errorw("message", "删除角色失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) FindAll(ctx context.Context, req *FindAllRequest) (*Roles, error) {
	roles, err := ctrl.repo.FindAll(ctx, req)
	if err != nil {
		ctrl.log.Errorw("message", "查询角色列表失败", "error", err.Error())
		return nil, err
	}

	return roles, nil
}

func (ctrl *Controller) FindRouteByRoles(ctx context.Context, req *FindRouteByRolesRequest) (*FindRouteByRolesResponse, error) {
	roles, err := ctrl.repo.FindByRoles(ctx, req.Roles)
	if err != nil {
		ctrl.log.Errorw("message", "查询角色列表失败", "error", err.Error())
		return nil, err
	}

	var routes []string
	for _, role := range roles {
		for _, menu := range role.Menus {
			routes = append(routes, menu.Router)
		}
	}

	// 去重
	removeDuplicate := slice.RemoveDuplicate(routes)

	return &FindRouteByRolesResponse{Routes: removeDuplicate}, nil
}
