package helpers

import (
	"fmt"
	"os"
	"time"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	auth_error "github.com/danzBraham/halo-suster/internal/exceptions/auth"
	"github.com/golang-jwt/jwt/v5"
)

var key = []byte(os.Getenv("JWT_SECRET"))

type CustomClaims struct {
	UserId string           `json:"userId"`
	Role   user_entity.Role `json:"role"`
	jwt.RegisteredClaims
}

func CreateJWT(ttl time.Duration, userId string, role user_entity.Role) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)

	claims := CustomClaims{
		userId,
		role,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

type Credential struct {
	UserId string           `json:"userId"`
	Role   user_entity.Role `json:"role"`
}

func VerifyJWT(tokenString string) (*Credential, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, auth_error.ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, auth_error.ErrUnknownClaims
	}

	return &Credential{
		UserId: claims.UserId,
		Role:   claims.Role,
	}, nil
}
