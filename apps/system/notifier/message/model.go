package message

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	// 消息Id
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 账户Id
	AccountId int64 `json:"accountId" bson:"account_id"`
	// 通道ApiKey
	ChannalApiKey string `json:"channalApiKey" bson:"channel_api_key"`
	// 状态
	Status string `json:"status" bson:"status"`
	// 消息体
	Body Body `json:"body" bson:"inline"`
	// 附加信息
	Payload string `json:"payload" bson:"payload"`
	// 创建时间
	CreatedAt int64 `json:"createdAt" bson:"created_at"`
	// 更新时间
	UpdatedAt int64 `json:"updatedAt" bson:"updated_at"`
}

// 消息体
type Body struct {
	// 等级
	Level string `json:"level" bson:"level"`
	// 类型
	Type string `json:"type" bson:"type"`
	// 标题
	Title string `json:"title" bson:"title"`
	// 正文
	Content string `json:"content" bson:"content"`
	// 加密正文, 仅在账户控制台可查看
	EncryptedContent string `json:"encryptedContent" bson:"encrypted_content"`
}

func (b *Body) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
