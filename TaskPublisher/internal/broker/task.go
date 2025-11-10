package broker

import (
	"context"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/dto"
	storageDTO "github.com/coffee-realist/TaskManager/TaskPublisher/internal/storage/dto"
)

type TaskBrokerInteractor interface {
	GetAllFinishedByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan storageDTO.TaskResp, error)
	Publish(task dto.TaskResp) error
}

type TaskBroker struct {
	interactor Interactor
}

func NewTaskBroker(interactor Interactor) *TaskBroker {
	return &TaskBroker{interactor: interactor}
}

func (t TaskBroker) GetAllFinishedByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan storageDTO.TaskResp, error) {
	tasks, err := t.interactor.Subscribe(ctx, taskReq.Project)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskBroker) Publish(task dto.TaskResp) error {
	err := t.interactor.Publish(task)
	if err != nil {
		return err
	}

	return nil
}
