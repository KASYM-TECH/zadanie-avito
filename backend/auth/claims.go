package auth

import (
	jwt "github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Username string
}
