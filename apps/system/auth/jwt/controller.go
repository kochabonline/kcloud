package jwt

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kit/auth/jwt"
	"github.com/kochabonline/kit/core/crypto/bcrypt"
	"github.com/kochabonline/kit/core/tools"
	"github.com/kochabonline/kit/errors"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	jwtRedis          *jwt.JwtRedis
	accountController *account.Controller
	log               log.Helper
}

func NewController(jwtRedis *jwt.JwtRedis, accountController *account.Controller, log log.Helper) *Controller {
	return &Controller{
		jwtRedis:          jwtRedis,
		accountController: accountController,
		log:               log,
	}
}

func (ctrl *Controller) Login(ctx context.Context, req *LoginRequest) (*Jwt, error) {
	account, err := ctrl.accountController.FindByUsername(ctx, req.Username)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return nil, err
	}
	if account.Id == 0 {
		return nil, common.ErrorAccountNotExist
	}

	if err := bcrypt.ComparePassword(account.Password, req.Password); err != nil {
		ctrl.log.Errorw("message", "密码错误", "error", err.Error())
		return nil, common.ErrorAccountOrPassword
	}

	claims := map[string]any{
		"id":       account.Id,
		"username": account.Username,
		"role":     account.Role,
	}

	atk, err := ctrl.jwtRedis.Generate(ctx, claims)
	if err != nil {
		ctrl.log.Errorw("message", "生成jwt失败", "error", err.Error())
		return nil, err
	}

	return &Jwt{AccessToken: atk.AccessToken, RefreshToken: atk.RefreshToken}, nil
}

func (ctrl *Controller) Logout(ctx context.Context, req *LogoutRequest) error {
	accountId, err := tools.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取账户Id失败", "error", err.Error())
		return err
	}

	if err := ctrl.jwtRedis.Delete(ctx, accountId); err != nil {
		ctrl.log.Errorw("message", "删除jwt失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) Refresh(ctx context.Context, req *RefreshRequest) (string, error) {
	newRefreshToken, err := ctrl.jwtRedis.Refresh(ctx, req.RefreshToken)
	if err != nil {
		ctrl.log.Errorw("message", "刷新jwt失败", "error", err.Error())
		return "", err
	}

	return newRefreshToken, nil
}

func (ctrl *Controller) Validate(c *gin.Context) (map[any]any, error) {
	result := make(map[any]any)
	ctx := c.Request.Context()

	// 获取Authorization Header
	authHeader := c.GetHeader(common.AuthorizationHeader)
	if authHeader == "" {
		ctrl.log.Errorw("message", "缺少Authorization Header")
		return nil, common.ErrorMissAuthorizationHeader
	}

	// 解析jwt
	claims, err := ctrl.jwtRedis.Parse(ctx, authHeader)
	if err != nil {
		ctrl.log.Errorw("message", "解析jwt失败", "error", err.Error())
		return nil, err
	}
	id, ok := claims["id"].(float64)
	if !ok {
		ctrl.log.Errorw("message", "解析jwt中的账户id失败")
		return nil, common.ErrUnauthorized
	}
	username, ok := claims["username"].(string)
	if !ok {
		ctrl.log.Errorw("message", "解析jwt中的用户名称失败")
		return nil, common.ErrUnauthorized
	}
	role, ok := claims["role"].(float64)
	if !ok {
		ctrl.log.Errorw("message", "解析jwt中的用户角色失败")
		return nil, common.ErrUnauthorized
	}
	// 获取语言
	lang := c.GetHeader("Language")

	result["id"] = int64(id)
	result["username"] = username
	result["role"] = int(role)
	result["lang"] = lang

	return result, nil
}

func (ctrl *Controller) Kick(ctx context.Context, req *KickRequest) error {
	// 获取账户角色
	role, err := tools.CtxValue[int](ctx, "role")
	if err != nil {
		ctrl.log.Errorw("message", "获取账户角色失败", "error", err.Error())
		return err
	}
	// 获取账户id
	account, err := ctrl.accountController.FindByUsername(ctx, req.Username)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return err
	}

	// 不能踢出比自己权限高的用户
	if role < int(account.Role) {
		return errors.Forbidden("您没有权限踢出该账户")
	}

	// 删除jwt
	if err := ctrl.jwtRedis.Delete(ctx, account.Id); err != nil {
		ctrl.log.Errorw("message", "删除jwt失败", "error", err.Error())
		return err
	}

	return nil
}
