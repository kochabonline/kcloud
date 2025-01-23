package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/transport/http/middleware"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("auth")
	{
		group.POST("login", h.Login)
		group.POST("logout", h.Logout)
		group.POST("refresh", h.Refresh)
		// 只有管理员可以踢人下线
		group.POST("kick", middleware.PermissionVPEWithConfig(
			middleware.PermissionVPEConfig{
				AllowedRoles: []string{common.RoleSuperAdmin, common.RoleAdmin},
			},
		), h.Kick)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 登录
// @description 登录
// @tags auth
// @accept json
// @produce json
// @param request body LoginRequest true "登录请求"
// @success 200 {object} Jwt
// @router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.Login(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 登出
// @description 登出
// @tags auth
// @accept json
// @produce json
// @param request body LogoutRequest true "登出请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Logout(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 刷新令牌
// @description 刷新令牌
// @tags auth
// @accept json
// @produce json
// @param request body RefreshRequest true "刷新请求"
// @param Authorization header string true "Authorization Token"
// @success 200 {string} string
// @router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.Refresh(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 踢人下线
// @description 踢人下线
// @tags auth
// @accept json
// @produce json
// @param request body KickRequest true "踢人请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /auth/kick [post]
func (h *Handler) Kick(c *gin.Context) {
	var req KickRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Kick(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}
