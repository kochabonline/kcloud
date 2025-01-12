package message

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	// 消息Id
	Id primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	// 账户Id
	AccountId int64 `json:"account_id"`
	// 通道ApiKey
	ChannalApiKey string `json:"channel_api_key"`
	// 状态
	Status string `json:"status"`
	// 消息体
	Body Body `bson:"inline"`
	// 附加信息
	Payload string `json:"payload"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at"`
}

// 消息体
type Body struct {
	// 等级
	Level string `json:"level"`
	// 类型
	Type string `json:"type"`
	// 标题
	Title string `json:"title"`
	// 正文
	Content string `json:"content"`
	// 加密正文, 仅在账户控制台可查看
	EncryptedContent string `json:"encrypted_content"`
}

func (b *Body) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
