package channal

import (
	"context"
)

type Interface interface {
	// 创建通道
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
	// 根据apiKey查找通道
	FindByApiKey(ctx context.Context, apiKey string) (*Channal, error)
	// 通道列表
	FindAll(ctx context.Context, req *FindAllRequest) (*Channels, error)
	// 删除通道
	Delete(ctx context.Context, req *DeleteRequest) error
}

// 创建通道请求
type CreateRequest struct {
	// 通道名称
	Name string `validate:"required" label:"通道名称"`
	// 通道描述
	Description string `validate:"omitempty" label:"通道描述"`
	// 供应商类型
	ProviderType string `validate:"required,oneof=email dingtalk lark telegram" label:"供应商类型"`
	// 供应商
	// email: {"from":"", "smtp_server":"", "smtp_port":0, "smtp_username":"", "smtp_password":""}
	// dingtalk: {"webhook":"", "secret":""}
	// lark: {"webhook":"", "secret":""}
	// telegram: {"token":"", "chat_id":""}
	Provider string `validate:"required" label:"供应商"`
	// 限速器类型
	LimiterType string `validate:"omitempty,oneof=sliding_window token_bucket" label:"限速器类型"`
	// 限速器
	// sliding_window: {"window":0, "limit":0}
	// token_bucket: {"capacity":0, "rate":0}
	Limiter string `validate:"omitempty" label:"限速器"`
	// 是否加签
	Sign bool `validate:"omitempty" label:"是否加签"`
}

// 创建通道响应
type CreateResponse struct {
	ApiKey string `label:"apikey"`
	Secret string ` label:"密钥"`
}

// 根据apiKey查找通道请求
type FindByApiKeyRequest struct {
	// ApiKey
	ApiKey string `form:"apiKey" validate:"required" label:"api密钥"`
}

// 通道列表请求
type FindAllRequest struct {
	// 页码
	Page int `form:"page" validate:"omitempty,gt=0" label:"页码"`
	// 每页数量
	Size int `form:"size" validate:"omitempty,gt=0" label:"每页数量"`
	// 供应商类型
	ProviderType string `form:"providerType" validate:"omitempty,oneof=email dingtalk telegram" label:"供应商类型"`
	// 通道状态
	Status Status `form:"status" validate:"omitempty,oneof=enabled disabled" label:"通道状态"`
	// 创建时间
	CreatedAt int64 `form:"createdAt" validate:"omitempty" label:"创建时间"`
	// 关键字
	Keyword string `form:"keyword" validate:"omitempty" label:"关键字"`
}

// 通道列表响应
type Channels struct {
	// 总数
	Total int64 `label:"总数"`
	// 通道列表
	Items []*Channal `label:"通道列表"`
}

// 删除通道请求
type DeleteRequest struct {
	// 通道Id
	Id string `uri:"id" validate:"required" label:"通道Id"`
}
