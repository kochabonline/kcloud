package message

import (
	"context"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/apps/system/notifier/channal"
	"github.com/kochabonline/kit/core/bot"
	"github.com/kochabonline/kit/core/bot/dingtalk"
	"github.com/kochabonline/kit/core/bot/email"
	"github.com/kochabonline/kit/core/bot/lark"
	"github.com/kochabonline/kit/core/bot/telegram"
	"github.com/kochabonline/kit/core/crypto/hmac"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo              *Repository
	channalController *channal.Controller
	log               log.Helper
}

func NewController(repo *Repository, channalController *channal.Controller, log log.Helper) *Controller {
	return &Controller{
		repo:              repo,
		channalController: channalController,
		log:               log,
	}
}

func (ctrl *Controller) Create(ctx context.Context, query *CreateQuery, req *CreateRequest) error {
	// 查询通道是否存在
	channal, err := ctrl.channalController.FindByApiKey(ctx, query.ApiKey)
	if err != nil {
		ctrl.log.Errorw("message", "查询通道失败", "error", err.Error())
		return err
	}
	if channal == nil {
		return common.ErrorChannelNotExist
	}

	// 验签
	if channal.Secret != "" {
		if err := hmac.Verify(channal.Secret, query.Signature, query.Timestamp); err != nil {
			ctrl.log.Errorw("message", "签名验证失败", "error", err.Error())
			return err
		}
	}

	// 构造消息
	now := time.Now().Unix()
	message := &Message{
		AccountId:     channal.AccountId,
		ChannalApiKey: query.ApiKey,
		Status:        Pending,
		CreatedAt:     now,
		UpdatedAt:     now,
		Body: Body{
			Level:            req.Level,
			Type:             req.Type,
			Title:            req.Title,
			Content:          req.Content,
			EncryptedContent: req.EncryptedContent,
		},
	}
	if req.Level == "" {
		message.Body.Level = Info
	}
	if req.Type == "" {
		message.Body.Type = Text
	}

	if err := ctrl.repo.Create(ctx, message); err != nil {
		ctrl.log.Errorw("message", "创建消息失败", "error", err.Error())
		return err
	}

	return nil
}

// 根据供应商类型生成对应消息体
func (ctrl *Controller) GenerateBody(ctx context.Context, ch *channal.Channal, msg *Message) bot.Sendable {
	if msg.Body.Title == "" {
		msg.Body.Title = DefaultTitle
	}

	switch ch.ProviderType {
	case channal.ProviderTypeEmail:
		if msg.Body.Type == Text {
			return email.NewMessage().With().Type(email.Text).Subject(msg.Body.Title).Body(msg.Body.Content).Message()
		}
		if msg.Body.Type == Markdown {
			return email.NewMessage().With().Type(email.Markdown).Subject(msg.Body.Title).Body(msg.Body.Content).Message()
		}
	case channal.ProviderTypeDingTalk:
		if msg.Body.Type == Text {
			return dingtalk.NewTextMessage().With().Content(msg.Body.Content).Message()
		}
		if msg.Body.Type == Markdown {
			return dingtalk.NewMarkdownMessage().With().Title(msg.Body.Title).Text(msg.Body.Content).Message()
		}
	case channal.ProviderTypeLark:
		var tag string
		if msg.Body.Type == Text {
			tag = lark.Text
		}
		if msg.Body.Type == Markdown {
			tag = lark.Markdown
		}
		return lark.NewCardMessage().With().CardHeader(msg.Body.Title).CardElement(tag, msg.Body.Content).Message()
	case channal.ProviderTypeTelegram:
		var mode string
		if msg.Body.Type == Text {
			mode = telegram.HTML
		}
		if msg.Body.Type == Markdown {
			mode = telegram.MarkdownV2
		}
		provider, _ := ch.Provider.(channal.ProviderTelegram)
		return telegram.NewSendMessage().With().ChatId(provider.ChatId).Text(msg.Body.Content).ParseMode(mode).Message()
	}
	return nil
}

func (ctrl *Controller) ChangeStatus(ctx context.Context, req *ChangeStatusRequest) error {
	if err := ctrl.repo.ChangeStatus(ctx, req); err != nil {
		ctrl.log.Errorw("message", "更新消息状态失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) FindAll(ctx context.Context, req *FindAllRequest) (*Messages, error) {
	messages, err := ctrl.repo.FindAll(ctx, req)
	if err != nil {
		ctrl.log.Errorw("message", "查询消息失败", "error", err.Error())
		return nil, err
	}

	return messages, nil
}

func (ctrl *Controller) Delete(ctx context.Context, req *DeleteRequest) error {
	return nil
}
