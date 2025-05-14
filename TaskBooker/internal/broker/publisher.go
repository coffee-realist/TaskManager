package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/storage/dto"
)

type Publisher interface {
	Publish(task dto.TaskResp) error
}

func (n *NatsBrokerService) Publish(task dto.TaskResp) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("finished.%s", task.Project)
	_, err = n.js.Publish(context.Background(), subject, data)
	return err
}
