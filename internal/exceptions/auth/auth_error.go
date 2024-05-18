package auth_error

import "errors"

var (
	ErrInvalidToken  = errors.New("Invalid token")
	ErrUnknownClaims = errors.New("Unknown claims type")
)
