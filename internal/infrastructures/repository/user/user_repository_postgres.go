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

func (r *UserRepositoryPostgres) CreateITUser(ctx context.Context, payload *user_entity.RegisterITUser) (userId string, err error) {
	userId = ulid.Make().String()
	query := "INSERT INTO users (id, nip, name, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.DB.Exec(ctx, query, userId, &payload.NIP, &payload.Name, &payload.Password, user_entity.IT)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (r *UserRepositoryPostgres) CreateNurseUser(ctx context.Context, payload *user_entity.RegisterNurseUser) (userId string, err error) {
	userId = ulid.Make().String()
	query := "INSERT INTO users (id, nip, name, card_image_url, role) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.DB.Exec(ctx, query, userId, &payload.NIP, &payload.Name, &payload.CardImageURL, user_entity.Nurse)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (r *UserRepositoryPostgres) GetUserByNIP(ctx context.Context, nip int) (user *user_entity.User, err error) {
	user = &user_entity.User{}
	query := "SELECT id, nip, name, password, role FROM users WHERE nip = $1"
	err = r.DB.QueryRow(ctx, query, nip).Scan(&user.ID, &user.NIP, &user.Name, &user.Password, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user_error.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPostgres) GetUsers(ctx context.Context, params *user_entity.UserQueryParams) ([]*user_entity.UserList, error) {
	query := "SELECT id, nip, name, created_at FROM users WHERE 1 = 1"
	args := []interface{}{}
	argID := 1

	if params.UserID != "" {
		query += " AND id = $" + strconv.Itoa(argID)
		args = append(args, params.UserID)
		argID++
	}

	if params.NIP != "" {
		query += ` AND nip::VARCHAR LIKE '%' || $` + strconv.Itoa(argID) + ` || '%'`
		args = append(args, params.NIP)
		argID++
	}

	if params.Name != "" {
		query += ` AND LOWER(name) LIKE '%' || $` + strconv.Itoa(argID) + ` || '%'`
		args = append(args, params.Name)
		argID++
	}

	switch params.Role {
	case "it":
		query += " AND role = 'it'"
	case "nurse":
		query += " AND role = 'nurse'"
	}

	switch params.CreatedAt {
	case "asc":
		query += " ORDER BY created_at ASC"
	case "desc":
		query += " ORDER BY created_at DESC"
	}

	query += " LIMIT $" + strconv.Itoa(argID) + " OFFSET $" + strconv.Itoa(argID+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*user_entity.UserList{}
	for rows.Next() {
		var user user_entity.UserList
		if err := rows.Scan(&user.ID, &user.NIP, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
