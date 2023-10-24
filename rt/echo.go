package rt

import (
	"context"
	"time"

	"github.com/dxbuild/kit/mw"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

func ProvideEcho(log zerolog.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(mw.NewLogger(log))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	return e
}

func InvokeServer(lc fx.Lifecycle, cfg *Config, log zerolog.Logger, e *echo.Echo) {
	port := cfg.Get("PORT", "8080")
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info().Msgf("starting server on :%s", port)
				err := e.Start(":" + port)
				if err != nil {
					log.Error().Err(err).Msg("server failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, cfg.GetDuration("SHUTDOWN_TIMEOUT", 10*time.Second))
			log.Info().Msg("stopping server")
			defer cancel()
			return e.Shutdown(ctx)
		},
	})
}
