package user_error

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrNIPAlreadyExists = errors.New("NIP already exists")
	ErrUserIsNotIT      = errors.New("user is not IT")
	ErrUserIsNotNurse   = errors.New("user is not a nurse")
	ErrInvalidPassword  = errors.New("invalid password")
)
