package storage

import (
	"TaskBooker/internal/storage/dto"
	"database/sql"
)

type UserStorageInteractor interface {
	GetByUsername(username string) (dto.UserResp, error)
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetByUsername(username string) (dto.UserResp, error) {
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"
	var user dto.UserResp
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		return dto.UserResp{}, err
	}
	return user, nil
}
