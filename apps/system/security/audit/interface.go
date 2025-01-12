package audit

import "context"

type Interface interface {
	// 账户最近登录详情
	LastLoginDetail(ctx context.Context, req *LastLoginDetailRequest) error
}

// 账户登录环境请求
type LastLoginDetailRequest struct {
	// 设备名称
	DeviceName string `json:"device_name" validate:"required" label:"设备名称"`
	// 设备id
	DeviceId string `json:"device_id" validate:"required" label:"设备id"`
	// 用户客户端
	UserAgent string `json:"user_agent" validate:"required" label:"ua"`
}
