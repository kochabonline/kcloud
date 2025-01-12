package channal

import (
	"fmt"
	"net/http"

	"github.com/kochabonline/kcloud/internal/util"
	"github.com/kochabonline/kit/core/bot"
	"github.com/kochabonline/kit/core/bot/dingtalk"
	"github.com/kochabonline/kit/core/bot/email"
	"github.com/kochabonline/kit/core/bot/lark"
	"github.com/kochabonline/kit/core/bot/telegram"
	"github.com/kochabonline/kit/core/rate"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channal struct {
	// 消息Id
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 账户Id
	AccountId int64 `json:"account_id"`
	// 通道名称
	Name string `json:"name"`
	// 通道描述
	Description string `json:"description"`
	// ApiKey
	ApiKey string `json:"api_key"`
	// 密钥
	Secret string `json:"secret"`
	// 供应商类型
	ProviderType string `json:"provider_type"`
	// 供应商
	Provider any `json:"provider"`
	// 限速器类型
	LimiterType string `json:"limiter_type"`
	// 限速器
	Limiter any `json:"limiter"`
	// 通道状态
	Status Status `json:"status"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at"`
}

// 邮件供应商
type ProviderEmail struct {
	// SMTP服务器
	SMTPServer string `json:"smtp_server" validate:"required"`
	// SMTP端口
	SMTPPort int `json:"smtp_port" validate:"required,gt=0"`
	// SMTP用户名
	SMTPUsername string `json:"smtp_username" validate:"required"`
	// SMTP密码
	SMTPPassword string `json:"smtp_password" validate:"required"`
}

// 钉钉供应商
type ProviderDingTalk struct {
	// 机器人Webhook
	Webhook string `json:"webhook" validate:"required"`
	// 机器人密钥
	Secret string `json:"secret" validate:"omitempty"`
}

// 飞书供应商
type ProviderLark struct {
	// 机器人Webhook
	Webhook string `json:"webhook" validate:"required"`
	// 机器人密钥
	Secret string `json:"secret" validate:"omitempty"`
}

// Telegram供应商
type ProviderTelegram struct {
	// 机器人Token
	Token string `json:"token" validate:"required"`
	// 机器人ChatId
	ChatId int64 `json:"chat_id" validate:"required"`
}

// 滑动窗口限速器
type LimiterSlidingWindow struct {
	// 限速器窗口大小
	Window int `json:"window" validate:"required,gt=0"`
	// 限速器窗口单位
	Limit int `json:"limit" validate:"required,gt=0"`
}

// 令牌桶限速器
type LimiterTokenBucket struct {
	// 令牌桶容量
	Capacity int `json:"capacity" validate:"required,gt=0"`
	// 令牌桶速率
	Rate int `json:"rate" validate:"required,gt=0"`
}

// 获取限速器
func (ch *Channal) RateLimiter(client *redis.Client) (rate.Limiter, error) {
	switch ch.LimiterType {
	case LimiterTypeSlidingWindow:
		var limiter LimiterSlidingWindow
		if err := util.BsonUnmarshal(ch.Limiter, &limiter); err != nil {
			return nil, err
		}
		return rate.NewSlidingWindowLimiter(client, fmt.Sprintf(SlidingWindowKey, ch.AccountId), limiter.Window, limiter.Limit), nil
	case LimiterTypeTokenBucket:
		var limiter LimiterTokenBucket
		if err := util.BsonUnmarshal(ch.Limiter, &limiter); err != nil {
			return nil, err
		}
		return rate.NewTokenBucketLimiter(client, fmt.Sprintf(TokenBucketKey, ch.AccountId), limiter.Capacity, limiter.Rate), nil
	default:
		return nil, nil
	}
}

// 获取供应商
func (ch *Channal) Bot(client *http.Client) (bot.Bot, error) {
	switch ch.ProviderType {
	case ProviderTypeEmail:
		var provider ProviderEmail
		if err := util.BsonUnmarshal(ch.Provider, &provider); err != nil {
			return nil, err
		}
		return email.New(email.SmtpPlainAuth{
			Host:     provider.SMTPServer,
			Port:     provider.SMTPPort,
			Username: provider.SMTPUsername,
			Password: provider.SMTPPassword,
		}), nil
	case ProviderTypeDingTalk:
		var provider ProviderDingTalk
		if err := util.BsonUnmarshal(ch.Provider, &provider); err != nil {
			return nil, err
		}
		return dingtalk.New(provider.Webhook, provider.Secret, dingtalk.WithClient(client)), nil
	case ProviderTypeLark:
		var provider ProviderLark
		if err := util.BsonUnmarshal(ch.Provider, &provider); err != nil {
			return nil, err
		}
		return lark.New(provider.Webhook, provider.Secret, lark.WithClient(client)), nil
	case ProviderTypeTelegram:
		var provider ProviderTelegram
		if err := util.BsonUnmarshal(ch.Provider, &provider); err != nil {
			return nil, err
		}
		return telegram.New(provider.Token, telegram.WithClient(client)), nil
	default:
		return nil, nil
	}
}
