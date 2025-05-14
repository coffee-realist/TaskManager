package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskBooker/internal/domain/dto"
	storageDTO "github.com/coffee-realist/TaskManager/TaskBooker/internal/storage/dto"
)

type TaskStorageInteractor interface {
	Insert(task dto.TaskResp, userID int) (int, error)
	UpdateStatus(taskPatch dto.TaskPatch) error
	Delete(id int) (storageDTO.TaskResp, error)
}

type TaskStorage struct {
	db *sql.DB
}

func NewTaskStorage(db *sql.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (s *TaskStorage) Insert(task dto.TaskResp, userID int) (int, error) {
	query := `INSERT INTO tasks (name, project, description, status, publisher_id, 
                  booked_by, booked_at, status_updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				RETURNING id`
	var taskID int
	err := s.db.QueryRow(query, task.Name, task.Project, task.Description,
		task.Status, task.PublisherID, userID).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	if taskID == 0 {
		return 0, errors.New("task is already booked")
	}
	return taskID, nil
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
        RETURNING id, name, project, description, status, publisher_id, 
                  booked_by, booked_at, status_updated_at`
	var task storageDTO.TaskResp
	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.Name, &task.Project, &task.Description, &task.Status,
		&task.PublisherID, &task.BookedBy, &task.BookedAt, &task.StatusUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storageDTO.TaskResp{}, fmt.Errorf("task with id %d not found", id)
		}
		return storageDTO.TaskResp{}, err
	}

	return task, nil
}
