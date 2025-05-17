package auth

import (
	"fmt"
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
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(auth.aud),
		jwt.WithIssuer(auth.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
