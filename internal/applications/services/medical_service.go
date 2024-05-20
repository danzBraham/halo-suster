package services

import (
	"context"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	medical_entity "github.com/danzBraham/halo-suster/internal/domains/entities/medicals"
	"github.com/danzBraham/halo-suster/internal/domains/repositories"
	medical_error "github.com/danzBraham/halo-suster/internal/exceptions/medicals"
)

type MedicalService struct {
	MedicalRepository repositories.MedicalRepository
}

func NewMedicalService(medicalRepository repositories.MedicalRepository) interfaces.MedicalService {
	return &MedicalService{MedicalRepository: medicalRepository}
}

func (s *MedicalService) CreatePatient(ctx context.Context, payload *medical_entity.AddMedicalPatient) error {
	isIdentityNumberExists, err := s.MedicalRepository.VerifyIdentityNumber(ctx, payload.IdentityNumber)
	if err != nil {
		return err
	}
	if isIdentityNumberExists {
		return medical_error.ErrIdentityNumberAlreadyExists
	}

	err = s.MedicalRepository.CreatePatient(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
