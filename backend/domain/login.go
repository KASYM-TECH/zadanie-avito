package domain

type LoginRequest struct {
	Username string `validate:"required" json:"username"`
}

type LoginResponse struct {
	RefreshToken string
	AccessToken  string
}
