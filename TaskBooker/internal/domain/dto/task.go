package dto

type TaskReq struct {
	Project string `json:"project" uri:"project" binding:"required" example:"project1"`
}

type TaskResp struct {
	Name        string `json:"name" example:"Task1"`
	Project     string `json:"project" example:"project1"`
	Description string `json:"description" example:"Описание задачи"`
	Status      string `json:"status" example:"CREATED"`
	PublisherID int    `json:"publisher_id" example:"6"`
}

type TaskPatch struct {
	ID     int    `json:"task_id" example:"1"`
	Status string `json:"status" example:"BOOKED"`
}

type TaskFinish struct {
	ID int `json:"task_id" example:"1"`
}
