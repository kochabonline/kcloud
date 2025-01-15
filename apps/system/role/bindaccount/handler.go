package bindaccount

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("role/bind/account")
	{
		group.POST("bind", h.Bind)
		group.POST("unbind", h.Unbind)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

// @title kcloud API
// @summary 绑定账户角色
// @description 绑定账户角色
// @tags role
// @accept json
// @produce json
// @param request body request true "绑定账户角色请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role/bind/account/bind [post]
func (h *Handler) Bind(c *gin.Context) {
	var req request
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
// @summary 解绑账户角色
// @description 解绑账户角色
// @tags role
// @accept json
// @produce json
// @param request body request true "解绑账户角色请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role/bind/account/unbind [post]
func (h *Handler) Unbind(c *gin.Context) {
	var req request
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
