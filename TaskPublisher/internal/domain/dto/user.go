package dto

type LoginReq struct {
	Username string `json:"username" example:"p1"`
	Password string `json:"password" example:"pass6"`
}
