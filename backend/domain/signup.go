package domain

type SignupRequest struct {
	Login    string `validate:"required" body:"login"`
	RoleId   int    `validate:"required" body:"role_id"`
	Password string `validate:"required,gte=5" body:"password"`
}
