package medical_entity

import "time"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type AddMedicalPatient struct {
	IdentityNumber int    `json:"identityNumber" validate:"required,identitynumber"`
	PhoneNumber    string `json:"phoneNumber" validate:"required,min=10,max=15,startswith=+62"`
	Name           string `json:"name" validate:"required,min=3,max=30"`
	BirthDate      string `json:"birthDate" validate:"required,datetime=2006-01-02"`
	Gender         Gender `json:"gender" validate:"required,oneof=male female"`
	CardImageURL   string `json:"identityCardScanImg" validate:"required,url"`
}

type MedicalPatient struct {
	IdentityNumber int       `json:"identityNumber"`
	PhoneNumber    string    `json:"phoneNumber"`
	Name           string    `json:"name"`
	BirthDate      time.Time `json:"birthDate"`
	Gender         Gender    `json:"gender"`
	CreatedAt      time.Time `json:"createdAt"`
}

type MedicalPatientParams struct {
	IdentityNumber string
	Limit          string
	Offset         string
	Name           string
	PhoneNumber    string
	CreatedAt      string
}

type AddMedicalRecord struct {
	IdentityNumber int    `json:"identityNumber" validate:"required,identitynumber"`
	Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `json:"medications" validate:"required,min=1,max=2000"`
	UserID         string `json:"userId" validate:"required"`
}

type IdentityDetail struct {
	IdentityNumber int       `json:"identityNumber"`
	PhoneNumber    string    `json:"phoneNumber"`
	Name           string    `json:"name"`
	BirthDate      time.Time `json:"birthDate"`
	Gender         Gender    `json:"gender"`
	CardImageURL   string    `json:"identityCardScanImg"`
}

type CreatedByDetail struct {
	NIP    int    `json:"nip"`
	Name   string `json:"name"`
	UserID string `json:"userId"`
}

type MedicalRecord struct {
	IdentityDetail  IdentityDetail  `json:"identityDetail"`
	Symptoms        string          `json:"symptoms"`
	Medications     string          `json:"medications"`
	CreatedAt       time.Time       `json:"createdAt"`
	CreatedByDetail CreatedByDetail `json:"createdBy"`
}

type MedicalRecordParams struct {
	IdentityNumber string
	UserID         string
	NIP            string
	Limit          string
	Offset         string
	CreatedAt      string
}
