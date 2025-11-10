package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/dto"
	"github.com/nats-io/nats.go"
	"log/slog"
)

type Subscriber interface {
	Subscribe(ctx context.Context, project string) (<-chan dto.TaskResp, error)
}

func (n *NatsBrokerService) Subscribe(ctx context.Context, project string) (<-chan dto.TaskResp, error) {
	subject := fmt.Sprintf("created.%s", project)
	out := make(chan dto.TaskResp, 100)

	sub, err := n.jsClient.Subscribe(
		subject,
		func(msg *nats.Msg) {
			var task dto.TaskResp
			if err := json.Unmarshal(msg.Data, &task); err != nil {
				slog.Error("unmarshal failed", "err", err)
				return
			}
			select {
			case out <- task:
			case <-ctx.Done():
			}
		},
		nats.DeliverAll(),
		nats.AckNone(),
		nats.ReplayInstant(),
	)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		if err := sub.Unsubscribe(); err != nil {
			slog.Error("Failed to unsubscribe", "err", err)
		} else {
			slog.Info("Ephemeral consumer unsubscribed", "subject", subject)
		}

		close(out)
	}()

	return out, nil
}
