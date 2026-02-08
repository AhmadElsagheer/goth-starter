package authfx

import (
	"github.com/ahmad/gother-example/server/modules/auth/repo"
	"github.com/ahmad/gother-example/server/modules/auth/service"

	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		service.New,
		repo.New,
	),
)
