package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	user_repository "github.com/danzBraham/halo-suster/internal/domains/repositories/users"
	user_error "github.com/danzBraham/halo-suster/internal/exceptions/users"
	"github.com/danzBraham/halo-suster/internal/helpers"
)

type UserService struct {
	UserRepository user_repository.UserRepository
}

func NewUserService(userRepository user_repository.UserRepository) interfaces.UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) CreateITUser(ctx context.Context, user *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error) {
	id, err := s.UserRepository.VerifyNIP(ctx, strconv.Itoa(user.NIP))
	if err != nil && !errors.Is(err, user_error.UserNotFoundError) {
		return nil, err
	}
	if id != "" {
		return nil, user_error.NIPAlreadyExistsError
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	id, err = s.UserRepository.CreateITUser(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateJWT(2*time.Hour, id)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoggedInUser{
		ID:          id,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) UserLogin(ctx context.Context, payload *user_entity.LoginUser) (*user_entity.LoggedInUser, error) {
	user, err := s.UserRepository.GetByNIP(ctx, strconv.Itoa(payload.NIP))
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateJWT(2*time.Hour, user.ID)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoggedInUser{
		ID:          user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}
