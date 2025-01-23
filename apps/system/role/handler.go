package role

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("role")
	{
		group.POST("", h.Create)
		group.PUT("", h.Update)
		group.DELETE(":id", h.Delete)
		group.GET("", h.List)
		group.GET("route", h.Route)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 创建角色
// @description 创建角色
// @tags role
// @accept json
// @produce json
// @param request body Request true "创建角色请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role [post]
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

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 修改角色
// @description 修改角色
// @tags role
// @accept json
// @produce json
// @param request body Request true "修改角色请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role [put]
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

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 删除角色
// @description 删除角色
// @tags role
// @accept json
// @produce json
// @param request path DeleteRequest true "删除角色请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	if err := h.controller.Delete(c.Request.Context(), &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 角色列表
// @description 角色列表
// @tags role
// @accept json
// @produce json
// @param query query FindAllRequest true "角色列表请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /role [get]
func (h *Handler) List(c *gin.Context) {
	var req FindAllRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	roles, err := h.controller.FindAll(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, roles)
}

// @title kcloud API
// @summary 根据角色列表获取角色列表路由
// @description 根据角色列表获取角色列表路由
// @tags role
// @accept json
// @produce json
// @param query query FindRouteByRolesRequest true "根据角色列表获取角色列表路由请求"
// @param Authorization header string true "Authorization Token"
// @success 200 {object} FindRouteByRolesResponse
// @router /role/route [get]
func (h *Handler) Route(c *gin.Context) {
	var req FindRouteByRolesRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	routes, err := h.controller.FindRouteByRoles(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, routes)
}
