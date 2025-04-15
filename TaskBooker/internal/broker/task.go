package broker

import (
	"TaskBooker/internal/domain/dto"
	storageDTO "TaskBooker/internal/storage/dto"
	"context"
)

type TaskBrokerInteractor interface {
	GetAllByProject(taskReq dto.TaskReq, userID int) (<-chan dto.TaskResp, error)
	PublishFinishedTask(task storageDTO.TaskResp) error
	DeleteTask(task dto.TaskResp) error
}

type TaskBroker struct {
	interactor Interactor
}

func NewTaskBroker(interactor Interactor) *TaskBroker {
	return &TaskBroker{interactor: interactor}
}

func (t TaskBroker) GetAllByProject(taskReq dto.TaskReq, userID int) (<-chan dto.TaskResp, error) {
	tasks, err := t.interactor.Subscribe(context.Background(), taskReq.Project, userID)
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
