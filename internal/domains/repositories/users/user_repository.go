package user_repository

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserRepository interface {
	CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (id string, err error)
	GetUserByNIP(ctx context.Context, nip int) (user *user_entity.User, err error)
}
