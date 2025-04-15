package dto

type UserResp struct {
	Username       string `db:"username"`
	HashedPassword string `db:"hashedPassword"`
	ID             int    `db:"id"`
}
