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

func (s *UserService) CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error) {
	if strconv.Itoa(payload.NIP)[:3] != "615" {
		return nil, user_error.ErrUserIsNotIT
	}

	user, err := s.UserRepository.GetUserByNIP(ctx, payload.NIP)
	if err != nil && !errors.Is(err, user_error.ErrUserNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, user_error.ErrNIPAlreadyExists
	}

	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}
	payload.Password = hashedPassword

	id, err := s.UserRepository.CreateITUser(ctx, payload)
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateJWT(2*time.Hour, id, user_entity.IT)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoggedInUser{
		UserID:      id,
		NIP:         payload.NIP,
		Name:        payload.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) CreateNurseUser(ctx context.Context, payload *user_entity.RegisterNurseUser) (*user_entity.LoggedInUser, error) {
	if strconv.Itoa(payload.NIP)[:3] != "303" {
		return nil, user_error.ErrUserIsNotNurse
	}

	user, err := s.UserRepository.GetUserByNIP(ctx, payload.NIP)
	if err != nil && !errors.Is(err, user_error.ErrUserNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, user_error.ErrNIPAlreadyExists
	}

	userId, err := s.UserRepository.CreateNurseUser(ctx, payload)
	if err != nil {
		return nil, err
	}

	accessToken, err := helpers.CreateJWT(2*time.Hour, userId, user_entity.Nurse)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoggedInUser{
		UserID:      userId,
		NIP:         payload.NIP,
		Name:        payload.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) UserLogin(ctx context.Context, payload *user_entity.LoginUser) (*user_entity.LoggedInUser, error) {
	user, err := s.UserRepository.GetUserByNIP(ctx, payload.NIP)
	if err != nil {
		return nil, err
	}

	isMatch, err := helpers.MatchPassword(user.Password, payload.Password)
	if err != nil {
		return nil, err
	}
	if !isMatch {
		return nil, user_error.ErrInvalidPassword
	}

	accessToken, err := helpers.CreateJWT(2*time.Hour, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoggedInUser{
		UserID:      user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, params *user_entity.UserQueryParams) ([]*user_entity.UserList, error) {
	users, err := s.UserRepository.GetUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateNurseUser(ctx context.Context, payload *user_entity.UpdateNurseUser) error {
	currentUser, err := s.UserRepository.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return err
	}

	isNIPExists, err := s.UserRepository.VerifyNIP(ctx, payload.NIP)
	if err != nil {
		return err
	}
	if isNIPExists && payload.NIP != currentUser.NIP {
		return user_error.ErrNIPAlreadyExists
	}

	err = s.UserRepository.UpdateNurseUser(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
