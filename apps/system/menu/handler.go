package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("menu")
	{
		group.POST("", h.Create)
		group.PUT("", h.Update)
		group.DELETE(":id", h.Delete)
		group.GET("", h.Tree)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 创建菜单
// @description 创建菜单
// @tags menu
// @accept json
// @produce json
// @param request body Request true "创建菜单请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /menu [post]
func (h *Handler) Create(c *gin.Context) {
	var req Request
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Create(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
}

// @title kcloud API
// @summary 修改菜单
// @description 修改菜单
// @tags menu
// @accept json
// @produce json
// @param request body Request true "修改菜单请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /menu [put]
func (h *Handler) Update(c *gin.Context) {
	var req Request
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Update(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
}

// @title kcloud API
// @summary 删除菜单
// @description 删除菜单
// @tags menu
// @accept json
// @produce json
// @param request path DeleteRequest true "删除菜单请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /menu/{id} [delete]
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
}

// @title kcloud API
// @summary 树形菜单
// @description 树形菜单
// @tags menu
// @accept json
// @produce json
// @param token header string true "Authorization Token"
// @success 200
// @router /menu [get]
func (h *Handler) Tree(c *gin.Context) {
	resp, err := h.controller.Tree(c.Request.Context())
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}
