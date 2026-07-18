package auth

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
