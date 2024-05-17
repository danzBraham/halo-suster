package user_error

import "errors"

var (
	UserNotFoundError     = errors.New("user not found")
	NIPAlreadyExistsError = errors.New("NIP already exists")
	ValidationFailedError = errors.New("validation failed")
)
