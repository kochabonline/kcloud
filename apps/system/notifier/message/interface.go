package message

import "context"

type Interface interface {
	// 创建消息
	Create(ctx context.Context, query *CreateQuery, req *CreateRequest) error
	// 更新消息状态
	ChangeStatus(ctx context.Context, req *ChangeStatusRequest) error
	// 查找所有消息
	FindAll(ctx context.Context, req *FindAllRequest) (*Messages, error)
	// 删除消息
	Delete(ctx context.Context, req *DeleteRequest) error
}

// 创建消息查询参数
type CreateQuery struct {
	ApiKey    string `form:"api_key" validate:"required" label:"api密钥"`
	Timestamp int64  `form:"timestamp" validate:"omitempty" label:"时间戳"`
	Signature string `form:"signature" validate:"omitempty" label:"签名"`
}

// 创建消息请求
type CreateRequest struct {
	// 消息等级
	Level string `json:"level" validate:"omitempty,oneof=info warning critical" label:"消息等级"`
	// 消息类型
	Type string `json:"type" validate:"omitempty,oneof=text markdown" label:"消息类型"`
	// 消息标题
	Title string `json:"title" validate:"omitempty" label:"消息标题"`
	// 消息内容
	Content string `json:"content" validate:"required" label:"消息内容"`
	// 加密消息内容, 仅在账户控制台可查看
	EncryptedContent string `json:"encrypted_content" validate:"omitempty" label:"加密消息内容"`
}

// 更新消息请求
type ChangeStatusRequest struct {
	// id
	Id string `form:"id" validate:"required" label:"id"`
	// 状态
	Status string `form:"status" validate:"omitempty,oneof=pending success failure" label:"状态"`
}

// 查找所有消息请求
type FindAllRequest struct {
	// 页码
	Page int `form:"page" validate:"omitempty" label:"页码"`
	// 每页数量
	Size int `form:"size" validate:"omitempty" label:"每页数量"`
	// 消息等级
	Level string `form:"level" validate:"omitempty" label:"消息等级"`
	// 消息类型
	CreatedAt int `form:"created_at" validate:"omitempty" label:"创建时间"`
	// 关键字
	Keyword string `form:"keyword" validate:"omitempty" label:"关键字"`
}

// 所有消息
type Messages struct {
	// 总数
	Total int64 `json:"total"`
	// 消息列表
	Items []*Message `json:"items"`
}

// 删除消息请求
type DeleteRequest struct {
	// 消息id
	Id string `uri:"id" validate:"required" label:"id"`
}
