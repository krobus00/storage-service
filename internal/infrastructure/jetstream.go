package infrastructure

import (
	"github.com/krobus00/storage-service/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func NewJetstreamClient() (*nats.Conn, nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(config.JetstreamHost(),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logrus.Warn("nats disconnect")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logrus.Warn("nats reconnecting")
		}))
	if err != nil {
		return nil, nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(
		nats.PublishAsyncMaxPending(config.JetstreamMaxPending()),
	)

	if err != nil {
		return nc, nil, err
	}

	return nc, js, nil
}
