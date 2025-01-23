package google

import "context"

type Interface interface {
	// 生成谷歌验证器密钥二维码URL
	Generate(ctx context.Context) (*GenerateResponse, error)
	// 验证谷歌验证码
	Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error)
}

type GenerateResponse struct {
	// 谷歌验证器密钥
	Secret string `json:"secret" label:"密钥"`
	// 二维码base64
	QrCode string `json:"qrCode" label:"二维码"`
}

type ValidateRequest struct {
	// 谷歌验证码
	Code string `json:"code" validate:"required" label:"谷歌验证码"`
}

type ValidateResponse struct {
	// 是否验证通过
	Ok bool `json:"ok" label:"是否验证通过"`
	// 安全码
	SecurityCode string `json:"securityCode" label:"安全码"`
}
