package user_repository_postgres

import (
	"context"
	"errors"
	"strconv"

	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	user_repository "github.com/danzBraham/halo-suster/internal/domains/repositories/users"
	user_error "github.com/danzBraham/halo-suster/internal/exceptions/users"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type UserRepositoryPostgres struct {
	DB *pgxpool.Pool
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) user_repository.UserRepository {
	return &UserRepositoryPostgres{DB: db}
}

func (r *UserRepositoryPostgres) CreateITUser(ctx context.Context, user *user_entity.RegisterITUser) (id string, err error) {
	id = ulid.Make().String()
	nip := strconv.Itoa(user.Nip)
	query := "INSERT INTO users (id, nip, name, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.DB.Exec(ctx, query, id, nip, user.Name, user.Password, user_entity.IT)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepositoryPostgres) VerifyNIP(ctx context.Context, nip string) (id string, err error) {
	query := "SELECT id FROM users WHERE nip = $1"
	err = r.DB.QueryRow(ctx, query, nip).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", user_error.UserNotFoundError
	}
	if err != nil {
		return "", err
	}
	return id, nil
}
