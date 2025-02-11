/*
Package cmd provides the entry point for the application.
*/
package cmd

import (
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kcloud/apps/system/auth/google"
	"github.com/kochabonline/kcloud/apps/system/auth/jwt"
	"github.com/kochabonline/kcloud/apps/system/menu"
	"github.com/kochabonline/kcloud/apps/system/notifier/channal"
	"github.com/kochabonline/kcloud/apps/system/notifier/message"
	"github.com/kochabonline/kcloud/apps/system/role"
	"github.com/kochabonline/kcloud/apps/system/role/bindaccount"
	"github.com/kochabonline/kcloud/apps/system/role/bindmenu"
	"github.com/kochabonline/kcloud/apps/system/security/audit"
	"github.com/kochabonline/kcloud/apps/system/security/captcha"
	"github.com/kochabonline/kcloud/apps/system/security/device"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kcloud/internal/server"
	"github.com/kochabonline/kit/app"
	"github.com/kochabonline/kit/log"
	"github.com/kochabonline/kit/log/zerolog"
	"github.com/kochabonline/kit/transport"
)

func newApp(
	config *config.Config,
	log log.Helper,
	// jwt中间件使用
	jwtController *jwt.Controller,
	// 消息队列处理
	messageQueue *message.Queue,
	// 下面都是注册路由的handler
	jwtHandler *jwt.Handler,
	googleHandler *google.Handler,
	accountHandler *account.Handler,
	roleHandler *role.Handler,
	bindaccountHandler *bindaccount.Handler,
	bindmenuHandler *bindmenu.Handler,
	channalHandler *channal.Handler,
	messageHandler *message.Handler,
	captchaHandler *captcha.Handler,
	deviceHandler *device.Handler,
	auditHandler *audit.Handler,
	menuHandler *menu.Handler,
) *app.App {
	httpServer := server.NewHttpServer(
		config, jwtController,
		jwtHandler,
		googleHandler,
		accountHandler,
		roleHandler,
		bindmenuHandler,
		channalHandler,
		messageHandler,
		captchaHandler,
		deviceHandler,
		auditHandler,
		menuHandler,
	)

	// 启动消息队列处理
	go messageQueue.Handle()

	return app.NewApp(
		[]transport.Server{httpServer},
		app.WithCloseFuncs(messageQueue.Close),
	)
}

func run() {
	console := log.NewHelper(zerolog.New().With().Caller().Logger())
	app, cleanup, err := initializeApp(config.Cfg, console)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
