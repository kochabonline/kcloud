package system

import (
	"github.com/google/wire"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kcloud/apps/system/auth/google"
	"github.com/kochabonline/kcloud/apps/system/auth/jwt"
	"github.com/kochabonline/kcloud/apps/system/menu"
	"github.com/kochabonline/kcloud/apps/system/notifier/channal"
	"github.com/kochabonline/kcloud/apps/system/notifier/message"
	"github.com/kochabonline/kcloud/apps/system/security/audit"
	"github.com/kochabonline/kcloud/apps/system/security/captcha"
	"github.com/kochabonline/kcloud/apps/system/security/device"
)

var ProviderSet = wire.NewSet(
	account.ProviderSet,
	google.ProviderSet,
	jwt.ProviderSet,
	channal.ProviderSet,
	message.ProviderSet,
	captcha.ProviderSet,
	device.ProviderSet,
	audit.ProviderSet,
	menu.ProviderSet,
)
