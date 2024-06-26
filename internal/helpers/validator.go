package helpers

import (
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewValidate() {
	validate.RegisterValidation("nip", validateUserNIP)
	validate.RegisterValidation("identitynumber", validateIdentityNumber)
	validate.RegisterValidation("imageurl", validateImageURL)
	validate.RegisterValidation("iso8601date", validateISO8601Date)
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

	if len(nipStr) < 13 || len(nipStr) > 15 {
		return false
	}

	// Check first digits
	if nipStr[:3] != "615" && nipStr[:3] != "303" {
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
	if len(randomDigits) < 3 || len(randomDigits) > 5 {
		return false
	}
	if _, err := strconv.Atoi(randomDigits); err != nil {
		return false
	}

	return true
}

func validateIdentityNumber(fl validator.FieldLevel) bool {
	identityNumber := fl.Field().Int()
	identityNumberStr := strconv.Itoa(int(identityNumber))
	return len(identityNumberStr) == 16
}

func validateImageURL(fl validator.FieldLevel) bool {
	u, err := url.ParseRequestURI(fl.Field().String())
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	ext := path.Ext(u.Path)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return false
	}

	return true
}

func validateISO8601Date(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}
