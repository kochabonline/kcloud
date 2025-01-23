package jwt

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kit/auth/jwt"
	"github.com/kochabonline/kit/core/crypto/bcrypt"
	"github.com/kochabonline/kit/core/util"
	"github.com/kochabonline/kit/errors"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	jwtRedis    *jwt.JwtRedis
	accountCtrl *account.Controller
	log         log.Helper
}

func NewController(jwtRedis *jwt.JwtRedis, accountCtrl *account.Controller, log log.Helper) *Controller {
	return &Controller{
		jwtRedis:    jwtRedis,
		accountCtrl: accountCtrl,
		log:         log,
	}
}

func (ctrl *Controller) Login(ctx context.Context, req *LoginRequest) (*Jwt, error) {
	account, err := ctrl.accountCtrl.FindByUsername(ctx, req.Username)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "username", req.Username, "error", err.Error())
		return nil, common.ErrorAccountOrPassword
	}

	if err := bcrypt.ComparePassword(account.Password, req.Password); err != nil {
		ctrl.log.Errorw("message", "密码错误", "error", err.Error())
		return nil, common.ErrorAccountOrPassword
	}

	roles := make([]string, len(account.Roles))
	for i, role := range account.Roles {
		roles[i] = role.Name
	}
	claims := map[string]any{
		"id":       account.Id,
		"username": account.Username,
		"roles":    roles,
	}

	atk, err := ctrl.jwtRedis.Generate(ctx, claims)
	if err != nil {
		ctrl.log.Errorw("message", "生成jwt失败", "error", err.Error())
		return nil, err
	}

	return &Jwt{AccessToken: atk.AccessToken, RefreshToken: atk.RefreshToken}, nil
}

func (ctrl *Controller) Logout(ctx context.Context, req *LogoutRequest) error {
	accountId, err := util.CtxValue[int64](ctx, "id")
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
	ctx := c.Request.Context()

	// 获取Authorization Header
	authHeader := c.GetHeader(common.AuthorizationHeader)
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		ctrl.log.Errorw("message", "缺少Authorization Header", "header", c.Request.Header)
		return nil, common.ErrorMissAuthorizationHeader
	}

	// 解析jwt
	claims, err := ctrl.jwtRedis.Parse(ctx, authHeader)
	if err != nil {
		ctrl.log.Errorw("message", "解析jwt失败", "claims", claims, "error", err.Error())
		return nil, errors.Unauthorized("%s", err.Error())
	}

	result := make(map[any]any)
	result["id"] = jwt.JwtMapClaimsParse[int64](claims, "id")
	result["username"] = jwt.JwtMapClaimsParse[string](claims, "username")
	result["roles"] = jwt.JwtMapClaimsParse[[]string](claims, "roles")
	result["lang"] = c.GetHeader(common.LanguageHeader)

	return result, nil
}

func (ctrl *Controller) Kick(ctx context.Context, req *KickRequest) error {
	// 获取账户id
	account, err := ctrl.accountCtrl.FindByUsername(ctx, req.Username)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return err
	}

	// 删除jwt
	if err := ctrl.jwtRedis.Delete(ctx, account.Id); err != nil {
		ctrl.log.Errorw("message", "删除jwt失败", "error", err.Error())
		return err
	}

	return nil
}
