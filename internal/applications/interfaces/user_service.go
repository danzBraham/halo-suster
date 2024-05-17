package interfaces

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserService interface {
	CreateITUser(ctx context.Context, user *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error)
}
