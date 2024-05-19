package interfaces

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserService interface {
	CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error)
	CreateNurseUser(ctx context.Context, payload *user_entity.RegisterNurseUser) (*user_entity.LoggedInUser, error)
	LoginUser(ctx context.Context, payload *user_entity.LoginUser) (*user_entity.LoggedInUser, error)
	GetUsers(ctx context.Context, params *user_entity.UserQueryParams) ([]*user_entity.UserList, error)
	UpdateNurseUser(ctx context.Context, payload *user_entity.UpdateNurseUser) error
	DeleteNurseUser(ctx context.Context, userId string) error
	GiveAccessNurseUser(ctx context.Context, payload *user_entity.GiveAccessNurseUser) error
}
