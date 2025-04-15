package broker

import (
	"TaskPublisher/internal/storage/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
	"strconv"
	"time"
)

type Subscriber interface {
	Subscribe(ctx context.Context, project string, userID int) (<-chan dto.TaskResp, error)
}

func (n *NatsBrokerService) Subscribe(ctx context.Context, project string, userID int) (<-chan dto.TaskResp, error) {
	subject := fmt.Sprintf("finished.%s", project)
	msgChan := make(chan dto.TaskResp, 100)

	cons, err := n.js.CreateOrUpdateConsumer(ctx, n.config.StreamName, jetstream.ConsumerConfig{
		Name:          strconv.Itoa(userID),
		Durable:       strconv.Itoa(userID),
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: subject,
	})
	if err != nil {
		return nil, err
	}

	go n.processMessages(ctx, cons, msgChan)
	return msgChan, nil
}

func (n *NatsBrokerService) processMessages(
	ctx context.Context,
	cons jetstream.Consumer,
	msgChan chan<- dto.TaskResp,
) {
	defer close(msgChan)
	batchSize := 100
	maxWait := 5 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := cons.Fetch(batchSize, jetstream.FetchMaxWait(maxWait))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, jetstream.ErrNoMessages) {
					continue
				}
				slog.Error("Failed to fetch messages", "error", err)
				return
			}

			for msg := range msgs.Messages() {
				var task dto.TaskResp
				if err := json.Unmarshal(msg.Data(), &task); err != nil {
					slog.Error("Failed to unmarshal task", "error", err)
					continue
				}

				select {
				case msgChan <- task:
					if err := msg.Ack(); err != nil {
						slog.Error("Failed to ack message", "error", err)
					}
				case <-time.After(1 * time.Second):
					slog.Warn("Message delivery timeout")
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
