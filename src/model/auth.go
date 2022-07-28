package model

type Auth struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}
