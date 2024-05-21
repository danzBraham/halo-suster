package medical_error

import "errors"

var (
	ErrIdentityNumberAlreadyExists = errors.New("identity number already exists")
	ErrIdentityNumberIsNotExists   = errors.New("identity number is not exists")
)
