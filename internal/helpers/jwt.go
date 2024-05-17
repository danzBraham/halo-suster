package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte(os.Getenv("JWT_SECRET"))

func CreateJWT(ttl time.Duration, subject string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiry),
		Subject:   subject,
	})
	return token.SignedString(key)
}
