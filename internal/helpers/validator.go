package helpers

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewValidate() {
	validate.RegisterValidation("nip", validateUserNIP)
}

func ValidatePayload(payload interface{}) error {
	if err := validate.Struct(payload); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func validateUserNIP(fl validator.FieldLevel) bool {
	nip := fl.Field().Int()
	nipStr := strconv.Itoa(int(nip))

	if len(nipStr) != 13 {
		return false
	}

	// Check first digits
	if nipStr[:3] != "615" {
		return false
	}

	// Check fourth digit
	if nipStr[3] != '1' && nipStr[3] != '2' {
		return false
	}

	// Check year
	currentYear := time.Now().Year()
	year := nipStr[4:8]
	if yearInt, err := strconv.Atoi(year); err != nil || yearInt < 2000 || yearInt > currentYear {
		return false
	}

	// Check month
	month := nipStr[8:10]
	if monthInt, err := strconv.Atoi(month); err != nil || monthInt < 1 || monthInt > 12 {
		return false
	}

	// Check random digits
	randomDigits := nipStr[10:]
	if _, err := strconv.Atoi(randomDigits); err != nil {
		return false
	}

	return true
}
