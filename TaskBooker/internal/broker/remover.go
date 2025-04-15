package broker

import (
	"TaskBooker/internal/domain/dto"
	"context"
	"fmt"
	"log/slog"
	"strconv"
)

type Remover interface {
	Remove(taskResp dto.TaskResp) error
}

func (n *NatsBrokerService) Remove(task dto.TaskResp) error {
	key := fmt.Sprintf("%s:%s", task.Project, task.Name)
	entry, err := n.tasksKV.Get(context.Background(), key)
	if err != nil {
		return fmt.Errorf("task index not found: %w", err)
	}

	seq, err := strconv.ParseUint(string(entry.Value()), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid sequence format: %w", err)
	}

	stream, err := n.js.Stream(context.Background(), n.config.StreamName)
	if err != nil {
		return err
	}

	if err = stream.DeleteMsg(context.Background(), seq); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if err = n.tasksKV.Delete(context.Background(), key); err != nil {
		slog.Error("Failed to delete task index", "error", err)
	}

	return nil
}
