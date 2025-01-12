package audit

import (
	"context"
	"net"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kit/app"
	"github.com/kochabonline/kit/core/tools"
	"github.com/kochabonline/kit/log"
	"github.com/oschwald/geoip2-golang"
)

var geo *geoip2.Reader

func init() {
	// 延迟加载GeoLite2数据库
	if !config.Cfg.Audit.GeoLite2.Enabled {
		return
	}
	var err error
	geo, err = geoip2.Open(config.Cfg.Audit.GeoLite2.Path)
	if err != nil {
		log.Fatal(err)
	}
	app.AddCloseFuncs(func() { geo.Close() })
}

type Controller struct {
	repo *Repository
	log  log.Helper
}

func NewController(repository *Repository, log log.Helper) *Controller {
	return &Controller{repo: repository, log: log}
}

// 账户最近登录详情
func (ctrl *Controller) LastLoginDetail(ctx context.Context, ip string, req *LastLoginDetailRequest) error {
	accountId, err := tools.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文中获取账户id失败", "err", err.Error())
		return err
	}
	accountUsername, err := tools.CtxValue[string](ctx, "username")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文中获取账户名称失败", "err", err.Error())
		return err
	}

	now := time.Now().Unix()
	// 根据ip获取地理位置
	audit := &Audit{
		Meta:               common.Meta{CreatedAt: now},
		AccountId:          accountId,
		AccountUsername:    accountUsername,
		LastLoginTimestamp: now,
		LastLoginIp:        ip,
		LastLoginDevice:    req.DeviceName,
		LastLoginUserAgent: req.UserAgent,
	}

	if config.Cfg.Audit.GeoLite2.Enabled {
		record, err := geo.City(net.ParseIP(ip))
		if err != nil {
			ctrl.log.Errorw("message", "获取ip地理位置失败", "err", err.Error())
			return err
		}
		audit.LastLoginLocation = record.City.Names["zh-CN"]
	}

	// 查询是否存在该账户的登录记录
	result, err := ctrl.repo.FindByAccountId(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户登录记录失败", "audit", audit, "err", err.Error())
		return err
	}
	if result == nil {
		// 创建新的登录记录
		if err := ctrl.repo.Create(ctx, audit); err != nil {
			ctrl.log.Errorw("message", "创建账户登录记录失败", "audit", audit, "err", err.Error())
			return err
		}
	} else {
		// 更新登录记录
		if err := ctrl.repo.Update(ctx, audit); err != nil {
			ctrl.log.Errorw("message", "更新账户登录记录失败", "audit", result, "err", err.Error())
			return err
		}
	}

	return nil
}
