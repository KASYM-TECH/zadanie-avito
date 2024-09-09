package jwt

import (
	"avito/auth"
	"avito/auth/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func NewRefreshToken(userID, roleId int) (string, bool) {
	tokenId := "refresh-" + uuid.New().String()
	return NewToken(userID, roleId, tokenId, config.RefreshTokenExpiration)
}

func NewAccessToken(userID, roleId int) (string, bool) {
	tokenId := "access-" + uuid.New().String()
	return NewToken(userID, roleId, tokenId, config.AccessTokenExpiration)
}

func NewToken(userID, roleId int, tokenId string, expiration time.Duration) (string, bool) {
	tokenExpirationTime := time.Now().Add(expiration)
	tokenClaims := &auth.Claims{
		RoleId: roleId,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenId,
			Subject:   strconv.Itoa(userID),
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		return "", false
	}

	return tokenString, true
}
