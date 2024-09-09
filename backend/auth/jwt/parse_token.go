package jwt

import (
	"avito/auth"
	"avito/auth/config"
	jwt "github.com/golang-jwt/jwt"
)

func ParseToken(tokenString string) (claims *auth.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
