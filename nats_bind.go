package kit

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

func NatsBind[IN, OUT any](log zerolog.Logger, nc *nats.Conn, subject string, fn func(in IN, out *OUT) error) (func() error, error) {
	sub, err := nc.Subscribe(subject, func(m *nats.Msg) {
		err := NatsHandle(nc, m, fn)
		if err != nil {
			log.Error().Err(err).Msg("error handling request")
		}
	})
	if err != nil {
		return nil, err
	}
	return sub.Unsubscribe, nil
}

func NatsHandle[IN, OUT any](nc *nats.Conn, m *nats.Msg, fn func(in IN, out *OUT) error) error {
	var in IN
	err := json.Unmarshal(m.Data, &in)
	if err != nil {
		return err
	}
	var out OUT
	err = fn(in, &out)
	if err != nil {
		return err
	}
	data, err := json.Marshal(out)
	if err != nil {
		return err
	}
	err = nc.Publish(m.Reply, data)
	if err != nil {
		return err
	}
	return nil
}
