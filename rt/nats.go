package rt

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func ProvideNatsConn(lc fx.Lifecycle, log zerolog.Logger, cfg *Config) (*nats.Conn, error) {
	url := cfg.Get("NATS_URL", nats.DefaultURL)
	nc, err := nats.Connect(url)
	if err != nil {
		log.Error().Err(err).Msg("error connecting to nats")
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("closing nats connection")
			nc.Drain()
			return nil
		},
	})
	return nc, nil
}
