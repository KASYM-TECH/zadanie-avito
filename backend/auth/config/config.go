package config

import (
	"os"
	"time"
)

const (
	RefreshTokenExpiration = 5 * 24 * time.Hour
	AccessTokenExpiration  = time.Hour
	PasswordCost           = 4
)

var (
	JwtKey = []byte(os.Getenv("JWT_KEY"))
)
