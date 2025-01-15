package bindmenu

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("role/menu/edit")
	{
		group.POST("", h.Edit)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 编辑角色菜单
// @description 编辑角色菜单
// @tags role
// @accept json
// @produce json
// @param request body Request true "编辑角色菜单请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role/menu/edit [post]
func (h *Handler) Edit(c *gin.Context) {
	var req Request
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Edit(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}
