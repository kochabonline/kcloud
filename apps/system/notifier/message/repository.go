package message

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kochabonline/kcloud/internal/util"
	"github.com/kochabonline/kit/log"
	kkafka "github.com/kochabonline/kit/store/kafka"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	coll   *mongo.Collection
	reader *kafka.Reader
	writer *kafka.Writer
	log    log.Helper
}

func NewRepository(client *mongo.Client, k *kkafka.Kafka, log log.Helper) *Repository {
	repo := &Repository{
		coll:   client.Database("kcloud").Collection("message"),
		reader: k.ConsumerGroup("message", "message"),
		writer: k.AsyncProducer("message"),
		log:    log,
	}
	repo.writerCompletion()
	return repo
}

// 异步消费失败直接保存到数据库记录为失败
func (repo *Repository) writerCompletion() {
	repo.writer.Completion = func(messages []kafka.Message, err error) {
		if err != nil {
			for _, msg := range messages {
				var m Message
				if err := json.Unmarshal(msg.Value, &m); err != nil {
					repo.log.Errorw("message", string(msg.Value), "error", err.Error())
					continue
				}
				m.Status = Failure
				m.UpdatedAt = time.Now().Unix()
				m.Payload = err.Error()
				if _, err := repo.coll.InsertOne(context.Background(), m); err != nil {
					repo.log.Errorw("message", string(msg.Value), "error", err.Error())
				}
			}
		}
	}
}

func (repo *Repository) Create(ctx context.Context, m *Message) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte("message"),
		Value: bytes,
	}

	return repo.writer.WriteMessages(ctx, msg)
}

func (repo *Repository) ChangeStatus(ctx context.Context, req *ChangeStatusRequest) error {
	now := time.Now().Unix()
	_, err := repo.coll.UpdateOne(ctx, bson.M{"_id": req.Id}, bson.M{"$set": bson.M{"status": req.Status, "updated_at": now}})
	return err
}

func (repo *Repository) FindAll(ctx context.Context, req *FindAllRequest) (*Messages, error) {
	var messages Messages

	// 根据请求参数构建查询条件
	// TODO: 自己只能查看自己创建的消息
	filter := bson.M{}
	if req.Level != "" {
		filter["level"] = req.Level
	}
	if req.CreatedAt != 0 {
		filter["created_at"] = bson.M{"$gte": req.CreatedAt}
	}
	if req.Keyword != "" {
		filter["name"] = bson.M{"$regex": req.Keyword}
	}

	total, err := repo.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	messages.Total = total

	offset, limit := util.Paginate(req.Page, req.Size)

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := repo.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &messages.Items); err != nil {
		return nil, err
	}

	return &messages, nil
}

func (repo *Repository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.coll.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}

// 插入一条成功消息到数据库
func (repo *Repository) InsertSuccess(ctx context.Context, m *Message) error {
	m.Status = Success
	m.UpdatedAt = time.Now().Unix()
	_, err := repo.coll.InsertOne(ctx, m)
	return err
}

// 插入一条失败消息到数据库
func (repo *Repository) InsertFailure(ctx context.Context, m *Message, err error) error {
	m.Status = Failure
	m.UpdatedAt = time.Now().Unix()
	m.Payload = err.Error()
	_, err = repo.coll.InsertOne(ctx, m)
	return err
}
