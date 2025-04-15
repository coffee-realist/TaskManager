package broker

import (
	"TaskPublisher/internal/domain/dto"
	"context"
	"encoding/json"
	"fmt"
)

type Publisher interface {
	Publish(task dto.TaskResp) error
}

func (n *NatsBrokerService) Publish(task dto.TaskResp) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("created.%s", task.Project)
	ack, err := n.js.Publish(context.Background(), subject, data)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", task.Project, task.Name)
	_, err = n.tasksKV.Put(context.Background(), key, []byte(fmt.Sprint(ack.Sequence)))
	return err
}
