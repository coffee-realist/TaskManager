package storage

import (
	"context"
	"database/sql"
	"errors"
	storageDTO "github.com/coffee-realist/TaskManager/TaskBooker/internal/storage/dto"
)

type TokenStorageInteractor interface {
	Insert(req storageDTO.TokenReq) error
	DeleteByUserID(userID int) error
	GetUserIDByRefreshToken(refreshToken string) (int, error)
}

type TokenStorage struct {
	db *sql.DB
}

func NewTokenStorage(db *sql.DB) *TokenStorage { return &TokenStorage{db: db} }

func (t *TokenStorage) Insert(req storageDTO.TokenReq) error {
	_, err := t.db.ExecContext(context.Background(),
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		req.UserID, req.Token, req.ExpiresAt)
	return err
}

func (t *TokenStorage) DeleteByUserID(userID int) error {
	res, err := t.db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", userID)

	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return errors.New("token delete failed")
	}
	return nil
}

func (t *TokenStorage) GetUserIDByRefreshToken(token string) (int, error) {
	var userID int
	err := t.db.QueryRowContext(context.Background(),
		"SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()", token).
		Scan(&userID)
	return userID, err
}
