package user_error

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrNIPAlreadyExists = errors.New("NIP already exists")
	ErrNotITUserNIP     = errors.New("NIP not starts with 615")
	ErrNotNurseUserNIP  = errors.New("NIP not starts with 303")
	ErrInvalidPassword  = errors.New("invalid password")
)
