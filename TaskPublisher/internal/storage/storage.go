package storage

import (
	"TaskPublisher/internal/domain/config"
	"database/sql"
	"fmt"
	"log"
)

func NewSqlConnection(cfg config.DataBaseConfig) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type Storage struct {
	UserStorage  UserStorageInteractor
	TaskStorage  TaskStorageInteractor
	TokenStorage TokenStorageInteractor
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:  NewUserStorage(db),
		TaskStorage:  NewTaskStorage(db),
		TokenStorage: NewTokenStorage(db),
	}
}
