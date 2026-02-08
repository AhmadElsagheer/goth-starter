package authfx

import (
	"{{AUTH_MODULE}}/repo"
	"{{AUTH_MODULE}}/service"

	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		service.New,
		repo.New,
	),
)
