package broker

import (
	"context"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/dto"
	storageDTO "github.com/coffee-realist/TaskManager/TaskBooker/internal/storage/dto"
)

type TaskBrokerInteractor interface {
	GetAllByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan dto.TaskResp, error)
	PublishFinishedTask(task storageDTO.TaskResp) error
	DeleteTask(task dto.TaskResp) error
}

type TaskBroker struct {
	interactor Interactor
}

func NewTaskBroker(interactor Interactor) *TaskBroker {
	return &TaskBroker{interactor: interactor}
}

func (t TaskBroker) GetAllByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan dto.TaskResp, error) {
	tasks, err := t.interactor.Subscribe(ctx, taskReq.Project)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskBroker) PublishFinishedTask(task storageDTO.TaskResp) error {
	err := t.interactor.Publish(task)
	if err != nil {
		return err
	}

	return nil
}

func (t TaskBroker) DeleteTask(task dto.TaskResp) error {
	err := t.interactor.Remove(task)
	if err != nil {
		return err
	}
	return nil
}
