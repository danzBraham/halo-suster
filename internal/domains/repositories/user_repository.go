package repositories

import (
	"context"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
)

type UserRepository interface {
	VerifyNIP(ctx context.Context, nip int) (bool, error)
	CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (userId string, err error)
	CreateNurseUser(ctx context.Context, payload *user_entity.RegisterNurseUser) (userId string, err error)
	GetUserByNIP(ctx context.Context, nip int) (user *user_entity.User, err error)
	GetUserByID(ctx context.Context, id string) (user *user_entity.User, err error)
	GetUsers(ctx context.Context, params *user_entity.UserQueryParams) ([]*user_entity.UserList, error)
	UpdateNurseUser(ctx context.Context, payload *user_entity.UpdateNurseUser) error
	DeleteNurseUser(ctx context.Context, userId string) error
	GiveAccessNurseUser(ctx context.Context, payload *user_entity.GiveAccessNurseUser) error
}
