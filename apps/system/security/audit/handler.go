package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller *Controller
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("security/audit")
	{
		group.POST("", h.LastLoginDetail)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller: controller}
}

// @title kcloud API
// @summary 账户最近登录详情
// @description 账户最近登录详情
// @tags security
// @accept json
// @produce json
// @param request body LastLoginDetailRequest true "账户登录环境请求"
// param x-real-ip header string true "用户真实ip"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /security/audit [post]
func (h *Handler) LastLoginDetail(c *gin.Context) {
	var req LastLoginDetailRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
	ip := c.GetHeader("x-real-ip")
	if err := h.controller.LastLoginDetail(c.Request.Context(), ip, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
	response.GinJSON(c, nil)
}
