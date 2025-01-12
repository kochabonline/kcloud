package account

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
	group := r.Group("account")
	{
		// 只有管理员可以创建账户
		group.POST("", middleware.PermissionVPEWithConfig(
			middleware.PermissionVPEConfig{
				AllowedRole: int(common.RoleAdmin),
			},
		), h.Create)
		group.PUT("change/password", h.ChangePassword)
		group.PUT("change/email", h.ChangeEmail)
		group.PUT("change/mobile", h.ChangeMobile)
		// 只有管理员可以查看账户列表
		group.GET("", middleware.PermissionVPEWithConfig(
			middleware.PermissionVPEConfig{
				AllowedRole: int(common.RoleAdmin),
			},
		), h.List)
		group.GET("detail", h.Detail)
		group.DELETE(":id", h.Delete)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 账户创建
// @description 账户创建
// @tags account
// @accept json
// @produce json
// @param request body CreateRequest true "账户创建请求"
// @success 200
// @router /account [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
	
	if err := h.controller.Create(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 账户修改密码
// @description 账户修改密码
// @tags account
// @accept json
// @produce json
// @param request body ChangePasswordRequest true "账户修改密码请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /account/change/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.ChangePassword(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 账户修改邮箱
// @description 账户修改邮箱
// @tags account
// @accept json
// @produce json
// @param request body ChangeEmailRequest true "账户修改邮箱请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /account/change/email [put]
func (h *Handler) ChangeEmail(c *gin.Context) {
	var req ChangeEmailRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.ChangeEmail(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 账户修改手机号
// @description 账户修改手机号
// @tags account
// @accept json
// @produce json
// @param request body ChangeMobileRequest true "账户修改手机号请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /account/change/mobile [put]
func (h *Handler) ChangeMobile(c *gin.Context) {
	var req ChangeMobileRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.ChangeMobile(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 账户列表
// @description 获取账户列表
// @tags account
// @accept json
// @produce json
// @param query query FindAllRequest true "查询账户列表请求"
// @param token header string true "Authorization Token"
// @success 200 {object} Accounts
// @router /account [get]
func (h *Handler) List(c *gin.Context) {
	var req FindAllRequest
	if err := validator.GinShouldBindQuery(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.FindAll(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 账户详情
// @description 获取账户详情
// @tags account
// @accept json
// @produce json
// @param token header string true "Authorization Token"
// @success 200 {object} Account
// @router /account/detail [get]
func (h *Handler) Detail(c *gin.Context) {
	resp, err := h.controller.Detail(c.Request.Context())
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 账户删除
// @description 账户删除
// @tags account
// @accept json
// @produce json
// @param uri path DeleteRequest true "删除账户请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /account/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := validator.GinShouldBindUri(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Delete(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}
