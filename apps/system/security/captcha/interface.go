package captcha

import (
	"context"
)

type Interface interface {
	// 发送手机验证码
	SendMobileCode(ctx context.Context) error
	// 验证手机验证码, 返回剩余时间
	VerifyMobileCode(ctx context.Context, req *VerifyMobileCodeRequest) (*VerifyMobileCodeResponse, error)
	// 发送邮箱验证码
	SendEmailCode(ctx context.Context) error
	// 验证邮箱验证码, 返回剩余时间
	VerifyEmailCode(ctx context.Context, req *VerifyEmailCodeRequest) (*VerifyEmailCodeResponse, error)
}

// 验证手机验证码请求
type VerifyMobileCodeRequest struct {
	// 验证码
	Code string `json:"code" validate:"required"`
}

// 验证手机验证码响应
type VerifyMobileCodeResponse struct {
	// 是否验证成功
	Ok bool `json:"ok"`
	// 剩余时间
	Ttl int64 `json:"ttl"`
}

// 验证邮箱验证码请求
type VerifyEmailCodeRequest struct {
	// 验证码
	Code string `json:"code" validate:"required"`
}

// 验证邮箱验证码响应
type VerifyEmailCodeResponse struct {
	// 是否验证成功
	Ok bool `json:"ok"`
	// 剩余时间
	Ttl int64 `json:"ttl"`
}
