package jwt

import (
	"avito/auth"
)

func GetClaimsByToken(token string) (*auth.Claims, bool) {
	claims, err := ParseToken(token)
	if err != nil {
		return nil, false
	}
	return claims, true
}
