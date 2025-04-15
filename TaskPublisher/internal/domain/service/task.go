package service

import (
	"TaskPublisher/internal/broker"
	"TaskPublisher/internal/domain/dto"
	storageDTO "TaskPublisher/internal/storage/dto"
)

type TaskInteractor interface {
	GetAllFinishedByProject(taskReq dto.TaskReq, userID int) (<-chan storageDTO.TaskResp, error)
	Publish(task dto.TaskResp) error
}

type TaskService struct {
	broker broker.TaskBrokerInteractor
}

func NewTaskService(broker broker.TaskBrokerInteractor) *TaskService {
	return &TaskService{broker: broker}
}

func (s *TaskService) GetAllFinishedByProject(taskReq dto.TaskReq, userID int) (<-chan storageDTO.TaskResp, error) {
	tasks, err := s.broker.GetAllFinishedByProject(taskReq, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) Publish(task dto.TaskResp) error {
	err := s.broker.Publish(task)
	return err
}
