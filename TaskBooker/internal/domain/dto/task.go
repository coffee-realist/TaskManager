package dto

type TaskReq struct {
	Project string `json:"project"`
}

type TaskResp struct {
	Name        string `json:"name"`
	Project     string `json:"project"`
	Description string `json:"description"`
	Status      string `json:"status"`
	PublisherID int    `json:"publisher"`
}

type TaskPatch struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type TaskFinish struct {
	ID int `json:"task_id"`
}
