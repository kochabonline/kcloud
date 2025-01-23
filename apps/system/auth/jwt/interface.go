package jwt

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	// 登录
	Login(ctx context.Context, req *LoginRequest) (*Jwt, error)
	// 登出
	Logout(ctx context.Context, req *LogoutRequest) error
	// 刷新
	Refresh(ctx context.Context, req *RefreshRequest) (string, error)
	// 验证
	Validate(c *gin.Context) (map[any]any, error)
	// 踢人下线
	Kick(ctx context.Context, req *KickRequest) error
}

// 登录请求
type LoginRequest struct {
	// 用户名
	Username string `json:"username" validate:"required" label:"用户名"`
	// 密码
	Password string `json:"password" validate:"required" label:"密码"`
}

// 登出请求
type LogoutRequest struct {
	// 刷新令牌
	RefreshToken string `json:"refreshToken" validate:"required" label:"刷新令牌"`
}

type RefreshRequest struct {
	// 刷新令牌
	RefreshToken string `json:"refreshToken" validate:"required" label:"刷新令牌"`
}

type ValidateRequest struct {
	// 访问令牌
	AccessToken string `json:"accessToken" validate:"required" label:"访问令牌"`
}

type KickRequest struct {
	// 用户名
	Username string `json:"username" validate:"required" label:"用户名"`
}
