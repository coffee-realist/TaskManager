package service

import (
	"context"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/broker"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/dto"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/storage"
)

type TaskInteractor interface {
	GetAllByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan dto.TaskResp, error)
	Book(task dto.TaskResp, userID int) (int, error)
	Finish(id int) error
}

type TaskService struct {
	storage storage.TaskStorageInteractor
	broker  broker.TaskBrokerInteractor
}

func NewTaskService(taskStorage storage.TaskStorageInteractor, broker broker.TaskBrokerInteractor) *TaskService {
	return &TaskService{storage: taskStorage, broker: broker}
}

func (s *TaskService) GetAllByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan dto.TaskResp, error) {
	tasks, err := s.broker.GetAllByProject(ctx, taskReq)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) Book(task dto.TaskResp, userID int) (int, error) {
	taskID, err := s.storage.Insert(task, userID)
	if err != nil {
		return 0, err
	}

	err = s.broker.DeleteTask(task)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func (s *TaskService) Finish(id int) error {
	task, err := s.storage.Delete(id)
	if err != nil {
		return err
	}

	err = s.broker.PublishFinishedTask(task)
	if err != nil {
		return err
	}

	return nil
}
