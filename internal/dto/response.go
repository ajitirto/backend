package dto

type LoginResponse struct {
	Token string `json:"token"`
}

type ProfileResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
