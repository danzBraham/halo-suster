package user_repository

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserRepository interface {
	CreateITUser(ctx context.Context, user *user_entity.RegisterITUser) (id string, err error)
	VerifyNIP(ctx context.Context, nip string) (id string, err error)
	GetByNIP(ctx context.Context, nip string) (user *user_entity.User, err error)
}
