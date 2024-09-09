package domain

type LoginRequest struct {
	Login    string `validate:"required" body:"login"`
	Password string `validate:"required" body:"password"`
}

type LoginResponse struct {
	RefreshToken string
	AccessToken  string
}
