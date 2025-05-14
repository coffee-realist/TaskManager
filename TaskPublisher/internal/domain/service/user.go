package service

import (
	"errors"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/dto"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserInteractor interface {
	Login(loginReq dto.LoginReq) (int, error)
}

type UserService struct {
	storage storage.UserStorageInteractor
}

func NewUserService(storage storage.UserStorageInteractor) *UserService {
	return &UserService{storage: storage}
}

func (s UserService) Login(loginReq dto.LoginReq) (int, error) {
	user, err := s.storage.GetByUsername(loginReq.Username)
	if err != nil {
		return -1, err
	}
	if loginReq.Username != user.Username {
		return -1, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginReq.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, errors.New("invalid username or password")
		}
		return -1, err
	}

	return user.ID, nil
}
