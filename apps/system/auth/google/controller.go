package google

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kit/auth/gotp"
	"github.com/kochabonline/kit/core/util"
	"github.com/kochabonline/kit/log"
)

const (
	gaIssuer = "kcloud"
)

var (
	gaKey = "ga:%d"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	ga   gotp.GoogleAuthenticatorer
	repo *Repository
	log  log.Helper
}

func NewController(config *config.Config, repo *Repository, log log.Helper) *Controller {
	controller := &Controller{
		ga:   gotp.GA,
		repo: repo,
		log:  log,
	}

	return controller
}

func (ctrl *Controller) Generate(ctx context.Context) (*GenerateResponse, error) {
	username, err := util.CtxValue[string](ctx, "username")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取用户名失败", "error", err.Error())
		return nil, err
	}

	secret := gotp.GA.GenerateSecret()
	qrcode, err := gotp.GA.GenerateQRCode(username, gaIssuer, secret)
	if err != nil {
		ctrl.log.Errorw("message", "生成二维码失败", "error", err.Error())
		return nil, err
	}

	// 保存到redis
	accountId, err := util.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取账户Id失败", "error", err.Error())
		return nil, err
	}
	if err := ctrl.repo.Set(ctx, fmt.Sprintf(gaKey, accountId), secret, time.Duration(10*time.Minute)); err != nil {
		ctrl.log.Errorw("message", "保存验证码失败", "error", err.Error())
		return nil, err
	}

	return &GenerateResponse{Secret: secret, QrCode: qrcode}, nil
}

func (ctrl *Controller) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	accountId, err := util.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取账户Id失败", "error", err.Error())
		return nil, err
	}

	var secret string
	secret, err = ctrl.repo.Get(ctx, fmt.Sprintf(gaKey, accountId))
	if err != nil {
		ctrl.log.Errorw("message", "获取验证码失败", "error", err.Error())
		return nil, err
	}
	// redis中没有谷歌验证器密钥，查询数据库
	if secret == "" {
		ga, err := ctrl.repo.FindByAccountId(ctx, accountId)
		if err != nil {
			ctrl.log.Errorw("message", "查询GoogleAuth失败", "error", err.Error())
			return nil, err
		}
		secret = ga.Secret
	}

	ok, err := ctrl.ga.ValidateCode(secret, req.Code)
	if err != nil {
		ctrl.log.Errorw("message", "验证验证码失败", "error", err.Error())
		return nil, err
	}

	// 验证成功
	var securityCode string
	if ok {
		// 删除redis中的验证码
		if err := ctrl.repo.Del(ctx, fmt.Sprintf(gaKey, accountId)); err != nil {
			ctrl.log.Errorw("message", "删除验证码失败", "error", err.Error())
			return nil, err
		}

		// 保存到数据库
		ga := GoogleAuth{
			AccountId: accountId,
			Secret:    secret,
			Status:    common.StatusNormal,
		}
		if err := ctrl.repo.Create(ctx, ga); err != nil {
			ctrl.log.Errorw("message", "保存GoogleAuth失败", "error", err.Error())
			return nil, err
		}

		// 生成安全码
		securityCode = uuid.New().String()
		if err := ctrl.repo.Set(ctx, fmt.Sprintf(common.SecurityCodeKey, accountId), securityCode, time.Duration(3*time.Minute)); err != nil {
			ctrl.log.Errorw("message", "保存安全码失败", "error", err.Error())
			return nil, err
		}
	}

	return &ValidateResponse{Ok: ok, SecurityCode: securityCode}, nil
}
