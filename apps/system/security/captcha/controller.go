package captcha

import (
	"context"
	"fmt"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kit/core/bot/email"
	"github.com/kochabonline/kit/core/tools"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	accountController *account.Controller
	repo              *Repository
	email             *email.Email
	log               log.Helper
	expiration        int64
}

func NewController(controller *account.Controller, repo *Repository, config *config.Config, log log.Helper) *Controller {
	return &Controller{
		accountController: controller,
		repo:              repo,
		email: email.New(email.SmtpPlainAuth{
			Identity: "",
			Username: config.Email.Username,
			Password: config.Email.Password,
			Host:     config.Email.Host,
			Port:     config.Email.Port,
		}),
		log:        log,
		expiration: 60 * 5,
	}
}

func (ctrl *Controller) getMobile(ctx context.Context) (string, error) {
	var mobile string
	accountId, err := tools.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文中获取账户id失败", "error", err.Error())
		return mobile, common.ErrUnauthorized
	}
	account, err := ctrl.accountController.FindById(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return mobile, common.ErrUnauthorized
	}
	mobile = account.Mobile
	return mobile, nil
}

func (ctrl *Controller) SendMobileCode(ctx context.Context) error {
	mobile, err := ctrl.getMobile(ctx)
	if err != nil {
		ctrl.log.Errorw("message", "获取手机号失败", "error", err.Error())
		return err
	}
	// TODO: 实现发送手机验证码
	ctrl.log.Infow("message", "发送手机验证码", "mobile", mobile)
	return nil
}

func (ctrl *Controller) VerifyMobileCode(ctx context.Context, req *VerifyMobileCodeRequest) (*VerifyMobileCodeResponse, error) {
	var resp VerifyMobileCodeResponse
	mobile, err := ctrl.getMobile(ctx)
	if err != nil {
		ctrl.log.Errorw("message", "获取手机号失败", "error", err.Error())
		return &resp, err
	}

	captcha, ttl, err := ctrl.repo.GetMobileCode(ctx, mobile)
	if err != nil {
		ctrl.log.Errorw("message", "获取验证码失败", "error", err.Error())
		return &resp, err
	}
	resp.Ttl = int64(ttl.Seconds())

	if captcha != req.Code {
		ctrl.log.Errorw("message", "验证码错误", "code", req.Code)
		return &resp, common.ErrInvalidCaptcha
	}
	resp.Ok = true

	return &resp, nil
}

func (ctrl *Controller) getEmail(ctx context.Context) (string, error) {
	var email string
	accountId, err := tools.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文中获取账户id失败", "error", err.Error())
		return email, common.ErrUnauthorized
	}
	account, err := ctrl.accountController.FindById(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return email, common.ErrUnauthorized
	}
	email = account.Email
	return email, nil
}

func (ctrl *Controller) SendEmailCode(ctx context.Context) error {
	accountEmail, err := ctrl.getEmail(ctx)
	if err != nil {
		ctrl.log.Errorw("message", "获取邮箱失败", "error", err.Error())
		return err
	}

	code := tools.GenerateRandomCode(6)
	resp, err := ctrl.email.Send(email.NewMessage().With().
		From("kcloud").
		To([]string{accountEmail}).
		Subject("验证码").
		Body(fmt.Sprintf("你的验证码是: %s", code)).
		Message())
	if err != nil {
		ctrl.log.Errorw("message", "发送邮件失败", "email", accountEmail, "response", resp, "error", err.Error())
		return err
	}

	if err := ctrl.repo.SetEmailCode(ctx, accountEmail, code, time.Duration(ctrl.expiration)*time.Second); err != nil {
		ctrl.log.Errorw("message", "保存验证码失败", "email", accountEmail, "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) VerifyEmailCode(ctx context.Context, req *VerifyEmailCodeRequest) (*VerifyEmailCodeResponse, error) {
	var resp VerifyEmailCodeResponse
	email, err := ctrl.getEmail(ctx)
	if err != nil {
		ctrl.log.Errorw("message", "获取邮箱失败", "error", err.Error())
		return &resp, err
	}

	captcha, ttl, err := ctrl.repo.GetEmailCode(ctx, email)
	if err != nil {
		ctrl.log.Errorw("message", "获取验证码失败", "error", err.Error())
		return &resp, err
	}
	resp.Ttl = int64(ttl.Seconds())

	if captcha != req.Code {
		ctrl.log.Errorw("message", "验证码错误", "code", req.Code)
		return &resp, common.ErrInvalidCaptcha
	}
	resp.Ok = true

	return &resp, nil
}
