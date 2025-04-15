package storage

import (
	"TaskBooker/internal/domain/dto"
	storageDTO "TaskBooker/internal/storage/dto"
	"database/sql"
	"errors"
	"fmt"
)

type TaskStorageInteractor interface {
	Insert(task dto.TaskResp, userID int) error
	UpdateStatus(taskPatch dto.TaskPatch) error
	Delete(id int) (storageDTO.TaskResp, error)
}

type TaskStorage struct {
	db *sql.DB
}

func NewTaskStorage(db *sql.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (s *TaskStorage) Insert(task dto.TaskResp, userID int) error {
	query := `INSERT INTO tasks (name, project, description, status, publisher, 
                  bookedBy, bookedAt, StatusUpdatedAt) 
				VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				ON CONFLICT (name, project) DO NOTHING
				RETURNING id`
	var taskID int
	err := s.db.QueryRow(query, task.Name, task.Project, task.Description,
		task.Status, task.PublisherID, userID).Scan(&taskID)
	if err != nil {
		return err
	}
	if taskID == 0 {
		return errors.New("task is already booked")
	}
	return nil
}

func (s *TaskStorage) UpdateStatus(taskPatch dto.TaskPatch) error {
	query := `UPDATE tasks SET status = $1, bookedAt = CURRENT_TIMESTAMP WHERE id = $2`
	row := s.db.QueryRow(query, taskPatch.Status, taskPatch.ID)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (s *TaskStorage) Delete(id int) (storageDTO.TaskResp, error) {
	query := `
        DELETE FROM tasks 
        WHERE id = $1
        RETURNING id, name, project, description, status, publisher, 
                  bookedBy, bookedAt, StatusUpdatedAt`
	var task storageDTO.TaskResp
	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.Name, &task.Project, &task.Description, &task.Status,
		&task.Publisher, &task.BookedBy, &task.BookedAt, &task.StatusUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storageDTO.TaskResp{}, fmt.Errorf("task with id %d not found", id)
		}
		return storageDTO.TaskResp{}, err
	}

	return task, nil
}
