package handler

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
