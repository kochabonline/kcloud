package message

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("notifier/message")
	{
		group.POST("send", h.Create)
		group.GET("", h.List)
		group.PUT("", h.ChangeStatus)
		group.DELETE(":id", h.Delete)
	}
}

// @title kcloud API
// @summary 发送消息
// @description 发送消息
// @tags notifier
// @accept json
// @produce json
// @param query query CreateQuery true "创建消息查询"
// @param request body CreateRequest true "创建消息请求"
// @success 200
// @router /notifier/message/send [post]
func (h *Handler) Create(c *gin.Context) {
	var query CreateQuery
	if err := validator.GinShouldBindQuery(c, &query); err != nil {
		response.GinJSONError(c, err)
		return
	}
	var req CreateRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	err := h.controller.Create(c.Request.Context(), &query, &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 消息列表
// @description 消息列表
// @tags notifier
// @accept json
// @produce json
// @param query query FindAllRequest true "查询消息列表请求"
// @param Authorization header string true "Authorization Token"
// @success 200 {object} Messages
// @router /notifier/message [get]
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
// @summary 更新消息状态
// @description 更新消息状态
// @tags notifier
// @accept json
// @produce json
// @param request body ChangeStatusRequest true "更新消息状态请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /notifier/message [put]
func (h *Handler) ChangeStatus(c *gin.Context) {
	var req ChangeStatusRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	err := h.controller.ChangeStatus(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 删除消息
// @description 删除消息
// @tags notifier
// @accept json
// @produce json
// @param uri path DeleteRequest true "删除消息请求"
// @param Authorization header string true "Authorization Token"
// @success 200
// @router /notifier/message/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := validator.GinShouldBindUri(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	err := h.controller.Delete(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}
