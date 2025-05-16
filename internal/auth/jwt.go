package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuth(secret string, aud string, iss string) *JWTAuth {
	return &JWTAuth{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (auth *JWTAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(auth.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (auth *JWTAuth) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}
