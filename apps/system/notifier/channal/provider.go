package channal

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHandler, NewController, NewRepository)
