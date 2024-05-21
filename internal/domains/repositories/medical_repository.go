package repositories

import (
	"context"

	medical_entity "github.com/danzBraham/halo-suster/internal/domains/entities/medicals"
)

type MedicalRepository interface {
	VerifyIdentityNumber(ctx context.Context, identityNumber int) (bool, error)
	CreatePatient(ctx context.Context, payload *medical_entity.AddMedicalPatient) error
	GetMedicalPatients(ctx context.Context, params *medical_entity.MedicalPatientParams) ([]*medical_entity.MedicalPatient, error)
	CreateMedicalRecord(ctx context.Context, payload *medical_entity.AddMedicalRecord) error
}
