package channal

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("notifier/channal")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("detail", h.FindByApiKey)
		group.DELETE(":id", h.Delete)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 创建通道
// @description 创建通道
// @tags notifier
// @accept json
// @produce json
// @param request body CreateRequest true "创建通道请求"
// @param token header string true "Authorization Token"
// @success 200 {object} CreateResponse
// @router /notifier/channal [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.Create(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 查询通道
// @description 查询通道
// @tags notifier
// @accept json
// @produce json
// @param query query FindAllRequest true "查询通道列表请求"
// @param token header string true "Authorization Token"
// @success 200 {object} Channels
// @router /notifier/channal [get]
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
// @summary 查询通道
// @description 查询通道
// @tags notifier
// @accept json
// @produce json
// @param query query FindByApiKeyRequest true "查询通道请求"
// @param token header string true "Authorization Token"
// @success 200 {object} Channal
// @router /notifier/channal/detail [get]
func (h *Handler) FindByApiKey(c *gin.Context) {
	var req FindByApiKeyRequest
	if err := validator.GinShouldBindQuery(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}
	resp, err := h.controller.FindByApiKey(c.Request.Context(), req.ApiKey)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 删除通道
// @description 删除通道
// @tags notifier
// @accept json
// @produce json
// @param uri path DeleteRequest true "删除通道请求"
// @param token header string true "Authorization Token"
// @success 200
// @router /notifier/channal/{id} [delete]
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
