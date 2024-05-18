package user_error

import "errors"

var (
	ErrUserNotFound     = errors.New("User not found")
	ErrNIPAlreadyExists = errors.New("NIP already exists")
	ErrNotITUserNIP     = errors.New("NIP not starts with 615")
	ErrInvalidPassword  = errors.New("Invalid password")
)
