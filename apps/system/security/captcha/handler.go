package captcha

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http/response"
	"github.com/kochabonline/kit/validator"
)

type Handler struct {
	controller *Controller
}

func (h *Handler) Register(r gin.IRouter) {
	group := r.Group("captcha")
	{
		group.POST("mobile/send", h.SendMobileCode)
		group.POST("mobile/verify", h.VerifyMobileCode)
		group.POST("email/send", h.SendEmailCode)
		group.POST("email/verify", h.VerifyEmailCode)
	}
}

func NewHandler(controller *Controller) *Handler {
	return &Handler{controller: controller}
}

// @title kcloud API
// @summary 发送手机验证码
// @description 发送手机验证码
// @tags security
// @accept json
// @produce json
// @param token header string true "Authorization Token"
// @success 200
// @router /security/captcha/mobile/send [post]
func (h *Handler) SendMobileCode(c *gin.Context) {
	if err := h.controller.SendMobileCode(c.Request.Context()); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 验证手机验证码
// @description 验证手机验证码
// @tags security
// @accept json
// @produce json
// @param request body VerifyMobileCodeRequest true "验证手机验证码请求"
// @param token header string true "Authorization Token"
// @success 200 {object} VerifyMobileCodeResponse
// @router /security/captcha/mobile/verify [post]
func (h *Handler) VerifyMobileCode(c *gin.Context) {
	var req VerifyMobileCodeRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.VerifyMobileCode(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}

// @title kcloud API
// @summary 发送邮箱验证码
// @description 发送邮箱验证码
// @tags security
// @accept json
// @produce json
// @param token header string true "Authorization Token"
// @success 200
// @router /security/captcha/email/send [post]
func (h *Handler) SendEmailCode(c *gin.Context) {
	if err := h.controller.SendEmailCode(c.Request.Context()); err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, nil)
}

// @title kcloud API
// @summary 验证邮箱验证码
// @description 验证邮箱验证码
// @tags security
// @accept json
// @produce json
// @param request body VerifyEmailCodeRequest true "验证邮箱验证码请求"
// @param token header string true "Authorization Token"
// @success 200 {object} VerifyEmailCodeResponse
// @router /security/captcha/email/verify [post]
func (h *Handler) VerifyEmailCode(c *gin.Context) {
	var req VerifyEmailCodeRequest
	if err := validator.GinShouldBind(c, &req); err != nil {
		response.GinJSONError(c, err)
		return
	}

	resp, err := h.controller.VerifyEmailCode(c.Request.Context(), &req)
	if err != nil {
		response.GinJSONError(c, err)
		return
	}

	response.GinJSON(c, resp)
}
