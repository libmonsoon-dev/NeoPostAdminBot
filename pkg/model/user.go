package model

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`

	IsAdmin bool `json:"is_admin"`
}
