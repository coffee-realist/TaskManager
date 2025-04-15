package dto

import "time"

type Status string

const (
	StatusCreated   Status = "created"
	StatusBooked    Status = "booked"
	StatusCompleted Status = "completed"
)

type TaskResp struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Project         string    `json:"project" db:"project"`
	Description     string    `json:"description" db:"description"`
	Status          Status    `json:"status" db:"status"`
	Publisher       int       `json:"publisher" db:"publisher"`
	BookedBy        int       `json:"booked_by" db:"booked_by"`
	BookedAt        time.Time `json:"booked_at" db:"booked_at"`
	StatusUpdatedAt time.Time `json:"status_updated_at" db:"status_updated_at"`
}
