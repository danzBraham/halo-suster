package interfaces

import (
	"context"

	medical_entity "github.com/danzBraham/halo-suster/internal/domains/entities/medicals"
)

type MedicalService interface {
	CreatePatient(ctx context.Context, payload *medical_entity.AddMedicalPatient) error
	GetMedicalPatients(ctx context.Context, params *medical_entity.MedicalPatientParams) ([]*medical_entity.MedicalPatient, error)
}
