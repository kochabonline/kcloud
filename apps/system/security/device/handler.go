package device

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
	group := r.Group("security/device")
	group.Use(middleware.PermissionVPEWithConfig(
		middleware.PermissionVPEConfig{
			AllowedRole: int(common.RoleAdmin),
		},
	))
	{
		group.POST("bind", h.Bind)
		group.POST("unbind", h.Unbind)
		group.PUT("change", h.Change)
		group.GET("", h.List)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller: controller}
}

// @title kcloud API
// @summary 设备绑定
// @description 设备绑定
// @tags security
// @accept json
// @produce json
// @param request body BindRequest true "设备绑定请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /security/device/bind [post]
func (h *Handler) Bind(c *gin.Context) {
	var req BindRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Bind(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 设备解绑
// @description 设备解绑
// @tags security
// @accept json
// @produce json
// @param request body UnbindRequest true "设备解绑请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /security/device/unbind [post]
func (h *Handler) Unbind(c *gin.Context) {
	var req UnbindRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Unbind(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 设备更换
// @description 设备更换
// @tags security
// @accept json
// @produce json
// @param request body ChangeRequest true "设备更换请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /security/device/change [put]
func (h *Handler) Change(c *gin.Context) {
	var req ChangeRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Change(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 设备列表
// @description 设备列表
// @tags security
// @accept json
// @produce json
// @param query query FindAllRequest true "设备列表请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /security/device [get]
func (h *Handler) List(c *gin.Context) {
	var req FindAllRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	devices, err := h.controller.FindAll(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, devices)
}
