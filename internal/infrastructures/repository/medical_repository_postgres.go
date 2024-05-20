package repository_postgres

import (
	"context"
	"errors"
	"strconv"

	medical_entity "github.com/danzBraham/halo-suster/internal/domains/entities/medicals"
	"github.com/danzBraham/halo-suster/internal/domains/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type MedicalRepositoryPostgres struct {
	DB *pgxpool.Pool
}

func NewMedicalRepositoryPostgres(db *pgxpool.Pool) repositories.MedicalRepository {
	return &MedicalRepositoryPostgres{DB: db}
}

func (r *MedicalRepositoryPostgres) VerifyIdentityNumber(ctx context.Context, identityNumber int) (bool, error) {
	var isIdentityNumberExists int
	query := "SELECT 1 FROM patients WHERE identity_number = $1"
	err := r.DB.QueryRow(ctx, query, strconv.Itoa(identityNumber)).Scan(&isIdentityNumberExists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MedicalRepositoryPostgres) CreatePatient(ctx context.Context, payload *medical_entity.AddMedicalPatient) error {
	id := ulid.Make().String()
	query := `INSERT INTO 
							patients (id, identity_number, phone_number, name, birth_date, gender, card_image_url)
							VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(ctx, query,
		id,
		strconv.Itoa(payload.IdentityNumber),
		&payload.PhoneNumber,
		&payload.Name,
		&payload.BirthDate,
		&payload.Gender,
		&payload.CardImageURL)

	if err != nil {
		return err
	}

	return nil
}
