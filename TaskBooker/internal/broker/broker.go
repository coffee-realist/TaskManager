package broker

import (
	"context"
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"io"
)

type Interactor interface {
	Publisher
	Subscriber
	Remover
	io.Closer
}

type NatsBrokerService struct {
	js       jetstream.JetStream
	jsClient nats.JetStreamContext
	natsConn *nats.Conn
	config   config.NatsConfig
	tasksKV  jetstream.KeyValue
}

func NewNatsBrokerService(
	js jetstream.JetStream,
	jsClient nats.JetStreamContext,
	natsConn *nats.Conn,
	cfg config.NatsConfig,
) (*NatsBrokerService, error) {
	if _, err := js.Stream(context.Background(), cfg.StreamName); err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}
	kv, err := js.KeyValue(context.Background(), cfg.KVBucket)
	if err != nil {
		return nil, fmt.Errorf("KV bucket not found: %w", err)
	}

	return &NatsBrokerService{
		js:       js,
		jsClient: jsClient,
		natsConn: natsConn,
		config:   cfg,
		tasksKV:  kv,
	}, nil
}

func (n *NatsBrokerService) Close() error {
	return n.natsConn.Drain()
}
