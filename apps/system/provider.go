package system

import (
	"github.com/google/wire"
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
)

var ProviderSet = wire.NewSet(
	jwt.ProviderSet,
	google.ProviderSet,
	account.ProviderSet,
	role.ProviderSet,
	bindaccount.ProviderSet,
	bindmenu.ProviderSet,
	captcha.ProviderSet,
	device.ProviderSet,
	audit.ProviderSet,
	menu.ProviderSet,
	channal.ProviderSet,
	message.ProviderSet,
)
