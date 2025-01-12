package google

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller Interface
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("/auth/gotp")
	{
		group.POST("generate", h.Generate)
		group.POST("verify", h.Validate)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

// @title kcloud API
// @summary 生成谷歌验证器密钥二维码base64
// @description 生成谷歌验证器密钥二维码base64
// @tags auth
// @accept json
// @produce json
// @param token header string true "Authorization Token"
// @success 200 {object} GenerateResponse
// @router /auth/gotp/generate [post]
func (h *Handler) Generate(c *gin.Context) {
	resp, err := h.controller.Generate(c.Request.Context())
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 验证谷歌验证器
// @description 验证谷歌验证器
// @tags auth
// @accept json
// @produce json
// @param request body ValidateRequest true "验证请求"
// @param token header string true "Authorization Token"
// @success 200 {object} ValidateResponse
// @router /auth/gotp/verify [post]
func (h *Handler) Validate(c *gin.Context) {
	var req ValidateRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.Validate(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}
