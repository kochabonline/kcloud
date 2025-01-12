package channal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/core/tools"
	"github.com/kochabonline/kit/errors"
	"github.com/kochabonline/kit/log"
	"github.com/kochabonline/kit/validator"
	"github.com/rs/xid"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo *Repository
	log  log.Helper
}

func NewController(repo *Repository, log log.Helper) *Controller {
	return &Controller{
		repo: repo,
		log:  log,
	}
}

func (ctrl *Controller) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	accountId, err := tools.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取账户id失败", "error", err.Error())
		return nil, err
	}
	lang, err := tools.CtxValue[string](ctx, "lang")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取语言失败", "error", err.Error())
		return nil, err
	}

	var channal Channal
	channal.AccountId = accountId
	channal.Name = req.Name
	channal.ApiKey = xid.New().String()
	channal.Description = req.Description
	channal.ProviderType = req.ProviderType
	channal.LimiterType = req.LimiterType
	now := time.Now().Unix()
	channal.CreatedAt = now
	channal.UpdatedAt = now

	// 供应商
	var provider any
	switch req.ProviderType {
	case ProviderTypeEmail:
		provider = &ProviderEmail{}
	case ProviderTypeDingTalk:
		provider = &ProviderDingTalk{}
	case ProviderTypeLark:
		provider = &ProviderLark{}
	case ProviderTypeTelegram:
		provider = &ProviderTelegram{}
	}
	// 反序列化
	if err := json.Unmarshal([]byte(req.Provider), provider); err != nil {
		ctrl.log.Errorw("message", "供应商序列化失败", "error", err.Error())
		return nil, err
	}
	// 验证参数
	if err := validator.StructTrans(provider, lang); err != nil {
		return nil, errors.BadRequest("%v", err)
	}
	channal.Provider = provider

	// 限速器
	var limiter any
	switch req.LimiterType {
	case LimiterTypeSlidingWindow:
		limiter = &LimiterSlidingWindow{}
	case LimiterTypeTokenBucket:
		limiter = &LimiterTokenBucket{}
	default:
		limiter = ""
	}
	// 反序列化
	if req.Limiter != "" {
		if err := json.Unmarshal([]byte(req.Limiter), limiter); err != nil {
			ctrl.log.Errorw("message", "限速器序列化失败", "error", err.Error())
			return nil, err
		}
		// 验证参数
		if err := validator.StructTrans(limiter, lang); err != nil {
			return nil, errors.BadRequest("%v", err)
		}
		channal.Limiter = limiter
	}

	// 是否加签
	if req.Sign {
		channal.Secret = xid.New().String()
	}

	// 创建前查询是否存在
	if _, err := ctrl.repo.FindByName(ctx, channal.Name); err == nil {
		ctrl.log.Errorw("message", "通道名称已存在", "name", channal.Name)
		return nil, common.ErrorChannelExists
	}

	// 创建通道
	if err := ctrl.repo.Create(ctx, &channal); err != nil {
		ctrl.log.Errorw("message", "创建通道失败", "error", err.Error())
		return nil, err
	}

	return &CreateResponse{ApiKey: channal.ApiKey, Secret: channal.Secret}, nil
}

func (ctrl *Controller) FindByApiKey(ctx context.Context, apiKey string) (*Channal, error) {
	channal, err := ctrl.repo.FindByApiKey(ctx, apiKey)
	if err != nil {
		ctrl.log.Errorw("message", "根据apiKey查找通道失败", "error", err.Error())
		return nil, err
	}

	return channal, nil
}

func (ctrl *Controller) FindAll(ctx context.Context, req *FindAllRequest) (*Channels, error) {
	channals, err := ctrl.repo.FindAll(ctx, req)
	if err != nil {
		ctrl.log.Errorw("message", "查询通道列表失败", "error", err.Error())
		return nil, err
	}

	return channals, nil
}

func (ctrl *Controller) Delete(ctx context.Context, req *DeleteRequest) error {
	if err := ctrl.repo.Delete(ctx, req.Id); err != nil {
		ctrl.log.Errorw("message", "删除通道失败", "error", err.Error())
		return err
	}

	return nil
}
