package jwt

import (
	"avito/auth"
	"avito/auth/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

func NewRefreshToken(userID, username string) (string, bool) {
	tokenId := "refresh-" + uuid.New().String()
	return NewToken(userID, username, tokenId, config.RefreshTokenExpiration)
}

func NewAccessToken(userID, username string) (string, bool) {
	tokenId := "access-" + uuid.New().String()
	return NewToken(userID, username, tokenId, config.AccessTokenExpiration)
}

func NewToken(userID, username, tokenId string, expiration time.Duration) (string, bool) {
	tokenExpirationTime := time.Now().Add(expiration)
	tokenClaims := &auth.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenId,
			Subject:   userID,
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		return "", false
	}

	return "bearer " + tokenString, true
}
