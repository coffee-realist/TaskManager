package service

import (
	"TaskPublisher/internal/broker"
	"TaskPublisher/internal/storage"
	"database/sql"
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
		Task:   NewTaskService(taskBroker),
		Token:  NewTokenService(newStorage.TokenStorage),
		Broker: taskBroker,
	}, nil
}
