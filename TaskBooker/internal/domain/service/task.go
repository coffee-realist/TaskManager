package service

import (
	"TaskBooker/internal/broker"
	"TaskBooker/internal/domain/dto"
	"TaskBooker/internal/storage"
)

type TaskInteractor interface {
	GetAllByProject(taskReq dto.TaskReq, userID int) (<-chan dto.TaskResp, error)
	Book(task dto.TaskResp, userID int) error
	Finish(id int) error
}

type TaskService struct {
	storage storage.TaskStorageInteractor
	broker  broker.TaskBrokerInteractor
}

func NewTaskService(taskStorage storage.TaskStorageInteractor, broker broker.TaskBrokerInteractor) *TaskService {
	return &TaskService{storage: taskStorage, broker: broker}
}

func (s *TaskService) GetAllByProject(taskReq dto.TaskReq, userID int) (<-chan dto.TaskResp, error) {
	tasks, err := s.broker.GetAllByProject(taskReq, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) Book(task dto.TaskResp, userID int) error {
	err := s.storage.Insert(task, userID)
	if err != nil {
		return err
	}

	err = s.broker.DeleteTask(task)
	if err != nil {
		return err
	}

	return nil
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
