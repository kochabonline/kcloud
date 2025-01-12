package captcha

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHandler, NewController, NewRepository)