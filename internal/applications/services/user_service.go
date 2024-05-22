package services

import (
	"context"
	"strconv"
	"time"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	"github.com/danzBraham/halo-suster/internal/domains/repositories"
	user_error "github.com/danzBraham/halo-suster/internal/exceptions/users"
	"github.com/danzBraham/halo-suster/internal/helpers"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) interfaces.UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (*user_entity.LoggedInUser, error) {
	isNIPExists, err := s.UserRepository.VerifyNIP(ctx, payload.NIP)
	if err != nil {
		return nil, err
	}
	if isNIPExists {
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
	isNIPExists, err := s.UserRepository.VerifyNIP(ctx, payload.NIP)
	if err != nil {
		return nil, err
	}
	if isNIPExists {
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

func (s *UserService) LoginUser(ctx context.Context, payload *user_entity.LoginUser) (*user_entity.LoggedInUser, error) {
	isNIPExists, err := s.UserRepository.VerifyNIP(ctx, payload.NIP)
	if err != nil {
		return nil, err
	}
	if !isNIPExists {
		return nil, user_error.ErrNIPIsNotExists
	}

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
	isNIPExists, err := s.UserRepository.VerifyNIP(ctx, payload.NIP)
	if err != nil {
		return err
	}
	if isNIPExists {
		return user_error.ErrNIPAlreadyExists
	}

	currentUser, err := s.UserRepository.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return err
	}

	if strconv.Itoa(currentUser.NIP)[:3] != "303" {
		return user_error.ErrUserIsNotNurse
	}

	err = s.UserRepository.UpdateNurseUser(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteNurseUser(ctx context.Context, userId string) error {
	user, err := s.UserRepository.GetUserByID(ctx, userId)
	if err != nil {
		return err
	}

	if strconv.Itoa(user.NIP)[:3] != "303" {
		return user_error.ErrUserIsNotNurse
	}

	err = s.UserRepository.DeleteNurseUser(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GiveAccessNurseUser(ctx context.Context, payload *user_entity.GiveAccessNurseUser) error {
	user, err := s.UserRepository.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return err
	}

	if strconv.Itoa(user.NIP)[:3] != "303" {
		return user_error.ErrUserIsNotNurse
	}

	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	payload.Password = hashedPassword

	err = s.UserRepository.GiveAccessNurseUser(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
