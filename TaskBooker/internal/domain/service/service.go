package service

import (
	"database/sql"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/broker"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/storage"
)

type Service struct {
	User   UserInteractor
	Task   TaskInteractor
	Token  TokenInteractor
	Broker broker.TaskBrokerInteractor
}

func NewService(db *sql.DB, newBroker broker.Interactor) (*Service, error) {
	newStorage := storage.NewStorage(db)
	taskBroker := broker.NewTaskBroker(newBroker)
	return &Service{
		User:   NewUserService(newStorage.UserStorage),
		Task:   NewTaskService(newStorage.TaskStorage, taskBroker),
		Token:  NewTokenService(newStorage.TokenStorage),
		Broker: taskBroker,
	}, nil
}
