package auth

import (
	"fmt"
	"gw-currency-wallet/internal/storages"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret  = []byte("SecretSalt")
	ExpireTime time.Duration
)

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func GenerateToken(user storages.User, expireTime time.Duration) (string, error) {
	claims := Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token is malformed: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
