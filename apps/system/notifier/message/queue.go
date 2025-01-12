package message

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/kochabonline/kit/errors"
	"github.com/kochabonline/kit/log"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Queue struct {
	ctrl        *Controller
	repo        *Repository
	redisClient *redis.Client
	httpClient  *http.Client
	ctx         context.Context
	log         log.Helper
}

func NewQueue(controller *Controller, repository *Repository, client *redis.Client, log log.Helper) *Queue {
	queue := &Queue{
		ctrl:        controller,
		repo:        repository,
		redisClient: client,
		log:         log,
		httpClient:  httpClient(),
		ctx:         context.Background(),
	}

	return queue
}

// Handle 消息队列处理
func (q *Queue) Handle() {
	pool, err := ants.NewPool(256)
	if err != nil {
		q.log.Fatalw("message", "协程池创建失败", "error", err.Error())
	}
	defer pool.Release()

	for {
		select {
		case <-q.ctx.Done():
			return
		default:
			// ReadMessage(q.ctx)会尝试从消息队列中读取消息。
			// 如果消息队列中没有消息，读取操作可能会阻塞一段时间，直到有新消息到达。
			message, err := q.repo.reader.ReadMessage(q.ctx)
			if err != nil {
				if err.Error() == "EOF" || err.Error() == "fetching message: EOF" {
					return
				}
				q.log.Errorw("message", "消息队列消费失败", "error", err.Error())
				continue
			}
			if err := pool.Submit(func() { q.handle(message) }); err != nil {
				q.log.Errorw("message", "协程池提交任务失败", "msg", string(message.Value), "error", err.Error())
				continue
			}
		}
	}
}

// handle 消息队列实际处理
func (q *Queue) handle(message kafka.Message) {
	var msg Message
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		q.log.Errorw("message", "消息反序列化失败", "error", err.Error())
		return
	}

	// 下面的成功与否都需要入库
	channal, err := q.ctrl.channalController.FindByApiKey(q.ctx, msg.ChannalApiKey)
	if err != nil {
		q.repo.InsertFailure(q.ctx, &msg, err)
		q.log.Errorw("message", "查询通道失败", "channal", channal, "error", err.Error())
		return
	}

	// 限流
	limiter, err := channal.RateLimiter(q.redisClient)
	if err != nil {
		q.repo.InsertFailure(q.ctx, &msg, err)
		q.log.Errorw("message", "获取限流器失败", "channal", channal, "error", err.Error())
		return
	}
	if limiter == nil || limiter.Allow() {
		// 获取发送消息的供应商
		bot, err := channal.Bot(q.httpClient)
		if err != nil {
			q.repo.InsertFailure(q.ctx, &msg, err)
			q.log.Errorw("message", "获取供应商失败", "channal", channal, "error", err.Error())
			return
		}
		resp, err := bot.Send(q.ctrl.GenerateBody(q.ctx, channal, &msg))
		if err != nil {
			q.repo.InsertFailure(q.ctx, &msg, err)
			q.log.Errorw("message", "发送消息失败", "channal", channal, "error", err.Error())
			return
		}
		if resp.StatusCode != http.StatusOK {
			q.repo.InsertFailure(q.ctx, &msg, err)
			q.log.Errorw("message", "发送消息失败", "channal", channal, "error", errors.BadRequest("%v", resp))
			return
		}
		// 发送成功入库
		q.repo.InsertSuccess(q.ctx, &msg)
		return
	}

	//  限流后重新入队
	if err := q.repo.writer.WriteMessages(q.ctx, message); err != nil {
		q.repo.InsertFailure(q.ctx, &msg, err)
		q.log.Errorw("message", "消息重新入队失败", "error", err.Error())
		return
	}
}

func (q *Queue) Close() {
	q.ctx.Done()
}

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}
