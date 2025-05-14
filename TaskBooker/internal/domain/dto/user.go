package dto

type LoginReq struct {
	Username string `json:"username" example:"w1"`
	Password string `json:"password" example:"pass1"`
}
