package message

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHandler, NewController, NewQueue, NewRepository)
