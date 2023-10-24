package rt

import (
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(ProvideConfig),
	fx.Provide(ProvideEcho),
	fx.Provide(ProvideLogger),
	fx.Provide(ProvideGorm),
	fx.Provide(ProvideNatsConn),
	fx.Invoke(InvokeServer),
)
