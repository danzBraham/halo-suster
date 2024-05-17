package interfaces

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserService interface {
	CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error)
	UserLogin(ctx context.Context, payload *user_entity.LoginUser) (*user_entity.LoggedInUser, error)
}
