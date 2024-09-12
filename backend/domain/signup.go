package domain

type SignupRequest struct {
	Username  string `validate:"required,lte=100" json:"username"`
	FirstName string `validate:"required,lte=100" json:"firstname"`
	LastName  string `validate:"required,lte=100" json:"lastname"`
}
