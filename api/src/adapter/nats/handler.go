package nats

import (
	"github.com/nats-io/nats.go"
	"ravxcheckout/src/adapter/config"
)

var natsConn *nats.Conn

func getNatsConn() (*nats.Conn, error) {
	var err error
	cfg := config.GetConfig()

	if natsConn == nil {
		natsConn, err = nats.Connect(cfg.NatsConnString)

		if err != nil {
			return nil, err
		}
	}

	return natsConn, nil
}
