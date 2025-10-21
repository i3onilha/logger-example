package controller

import (
	"go.uber.org/fx"
)

// FX Module for controllers
var Module = fx.Module("controllers",
	fx.Provide(NewUserController),
)
