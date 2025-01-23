package operate

import "go.mongodb.org/mongo-driver/bson/primitive"

type Operate struct {
	// 自增id
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 操作类型
	Type string `json:"type" bson:"type"`
	// 操作时间
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	// 操作人
	Operator string `json:"operator" bson:"operator"`
	// 操作对象
	Object string `json:"object" bson:"object"`
	// 操作结果
	Result string `json:"result" bson:"result"`
	// 操作详情
	Detail string `json:"detail" bson:"detail"`
}
