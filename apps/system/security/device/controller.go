package device

import (
	"context"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo *Repository
	log  log.Helper
}

func NewController(repo *Repository, log log.Helper) *Controller {
	return &Controller{repo: repo, log: log}
}

func (ctrl *Controller) Bind(ctx context.Context, req *BindRequest) error {
	device := &Device{
		Meta: common.Meta{
			CreatedAt: time.Now().Unix(),
		},
		AccountId: req.AccountId,
		DeviceId:  req.DeviceId,
		Name:      req.Name,
		Type:      req.Type,
	}
	if err := ctrl.repo.Create(ctx, device); err != nil {
		ctrl.log.Errorw("message", "设备绑定失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) Unbind(ctx context.Context, req *UnbindRequest) error {
	device := &Device{
		AccountId: req.AccountId,
		DeviceId:  req.DeviceId,
	}
	if err := ctrl.repo.Delete(ctx, device); err != nil {
		ctrl.log.Errorw("message", "设备解绑失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) Change(ctx context.Context, req *ChangeRequest) error {
	device := &Device{
		AccountId: req.AccountId,
		DeviceId:  req.DeviceId,
		Name:      req.Name,
		Type:      req.Type,
	}

	if err := ctrl.repo.Update(ctx, device); err != nil {
		ctrl.log.Errorw("message", "设备更换失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) FindAll(ctx context.Context, req *FindAllRequest) (*Devices, error) {
	devices, err := ctrl.repo.FindAll(ctx, req)
	if err != nil {
		ctrl.log.Errorw("message", "设备列表查询失败", "error", err.Error())
		return nil, err
	}

	return devices, nil
}

// 根据账户id查询设备
func (ctrl *Controller) FindByAccountId(ctx context.Context, accountId string) (*Device, error) {
	device, err := ctrl.repo.FindByAccountId(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "根据账户id查询设备失败", "error", err.Error())
		return &Device{}, err
	}

	return device, nil
}
