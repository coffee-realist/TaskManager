package broker

import (
	"TaskPublisher/internal/domain/dto"
	storageDTO "TaskPublisher/internal/storage/dto"
	"context"
)

type TaskBrokerInteractor interface {
	GetAllFinishedByProject(taskReq dto.TaskReq, userID int) (<-chan storageDTO.TaskResp, error)
	Publish(task dto.TaskResp) error
}

type TaskBroker struct {
	interactor Interactor
}

func NewTaskBroker(interactor Interactor) *TaskBroker {
	return &TaskBroker{interactor: interactor}
}

func (t TaskBroker) GetAllFinishedByProject(taskReq dto.TaskReq, userID int) (<-chan storageDTO.TaskResp, error) {
	tasks, err := t.interactor.Subscribe(context.Background(), taskReq.Project, userID)
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
