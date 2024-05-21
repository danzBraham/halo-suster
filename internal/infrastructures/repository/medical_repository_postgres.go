package repository_postgres

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

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

func (r *MedicalRepositoryPostgres) GetMedicalPatients(ctx context.Context, params *medical_entity.MedicalPatientParams) ([]*medical_entity.MedicalPatient, error) {
	query := `SELECT identity_number, phone_number, name, birth_date, gender, created_at 
							FROM patients WHERE is_deleted = false`
	args := []interface{}{}
	argID := 1

	if params.IdentityNumber != "" {
		query += ` AND identity_number LIKE '%' || $` + strconv.Itoa(argID) + ` || '%'`
		args = append(args, params.IdentityNumber)
		argID++
	}

	if params.PhoneNumber != "" {
		query += ` AND phone_number LIKE '%' || $` + strconv.Itoa(argID) + ` || '%'`
		args = append(args, strings.TrimPrefix(params.PhoneNumber, "+"))
		argID++
	}

	if params.Name != "" {
		query += ` AND LOWER(name) LIKE '%' || $` + strconv.Itoa(argID) + ` || '%'`
		args = append(args, params.Name)
		argID++
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

	medicalPatients := []*medical_entity.MedicalPatient{}
	for rows.Next() {
		var identityNumber string
		var medicalPatient medical_entity.MedicalPatient
		err := rows.Scan(
			&identityNumber,
			&medicalPatient.PhoneNumber,
			&medicalPatient.Name,
			&medicalPatient.BirthDate,
			&medicalPatient.Gender,
			&medicalPatient.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		identityNumberInt, err := strconv.Atoi(identityNumber)
		if err != nil {
			return nil, err
		}
		medicalPatient.IdentityNumber = identityNumberInt
		medicalPatients = append(medicalPatients, &medicalPatient)
	}

	return medicalPatients, nil
}

func (r *MedicalRepositoryPostgres) CreateMedicalRecord(ctx context.Context, payload *medical_entity.AddMedicalRecord) error {
	id := ulid.Make().String()
	log.Println(payload.UserID)
	query := `INSERT INTO 
							medical_records (id, symptoms, medications, patient_identity_number, created_by)
							VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(ctx, query,
		id,
		&payload.Symptoms,
		&payload.Medications,
		strconv.Itoa(payload.IdentityNumber),
		&payload.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}
