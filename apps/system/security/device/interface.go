package device

import "context"

type Interface interface {
	// 设备列表
	FindAll(ctx context.Context, req *FindAllRequest) (*Devices, error)
	// 绑定设备
	Bind(ctx context.Context, req *BindRequest) error
	// 解绑设备
	Unbind(ctx context.Context, req *UnbindRequest) error
	// 更换设备
	Change(ctx context.Context, req *ChangeRequest) error
}

// 设备列表请求
type FindAllRequest struct {
	// 页码
	Page int `form:"page" validate:"omitempty"`
	// 每页数量
	Size int `form:"size" validate:"omitempty"`
	// 账户id
	AccountId string `form:"accountId" validate:"omitempty"`
	// 设备id
	DeviceId string `form:"deviceId" validate:"omitempty"`
}

// 设备列表响应
type Devices struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Items []*Device `json:"items"`
}

// 绑定设备请求
type BindRequest struct {
	// 账户id
	AccountId string `json:"accountId" validate:"required"`
	// 设备id
	DeviceId string `json:"deviceId" validate:"required"`
	// 设备名称
	Name string `json:"name" validate:"required"`
	// 设备类型
	Type string `json:"type" validate:"required"`
}

// 解绑设备请求
type UnbindRequest struct {
	// 账户id
	AccountId string `json:"accountId" validate:"required"`
	// 设备id
	DeviceId string `json:"deviceId" validate:"required"`
}

// 更换设备请求
type ChangeRequest struct {
	// 账户id
	AccountId string `json:"accountId" validate:"required"`
	// 设备id
	DeviceId string `json:"deviceId" validate:"required"`
	// 设备名称
	Name string `json:"name" validate:"required"`
	// 设备类型
	Type string `json:"type" validate:"required"`
}
