package config

import (
	"go.uber.org/fx"
)

// FX Module for config
var Module = fx.Module("config",
	fx.Provide(GetInstance),
)
