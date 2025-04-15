package dto

type TokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
