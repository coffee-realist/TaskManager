package service

import (
	"context"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/broker"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/dto"
	storageDTO "github.com/coffee-realist/TaskManager/TaskPublisher/internal/storage/dto"
)

type TaskInteractor interface {
	GetAllFinishedByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan storageDTO.TaskResp, error)
	Publish(task dto.TaskResp) error
}

type TaskService struct {
	broker broker.TaskBrokerInteractor
}

func NewTaskService(broker broker.TaskBrokerInteractor) *TaskService {
	return &TaskService{broker: broker}
}

func (s *TaskService) GetAllFinishedByProject(ctx context.Context, taskReq dto.TaskReq) (<-chan storageDTO.TaskResp, error) {
	tasks, err := s.broker.GetAllFinishedByProject(ctx, taskReq)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) Publish(task dto.TaskResp) error {
	err := s.broker.Publish(task)
	return err
}
