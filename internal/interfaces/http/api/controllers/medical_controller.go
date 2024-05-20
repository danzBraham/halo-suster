package controllers

import (
	"errors"
	"net/http"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	medical_entity "github.com/danzBraham/halo-suster/internal/domains/entities/medicals"
	medical_error "github.com/danzBraham/halo-suster/internal/exceptions/medicals"
	"github.com/danzBraham/halo-suster/internal/helpers"
	"github.com/danzBraham/halo-suster/internal/interfaces/http/api/middlewares"
	"github.com/go-chi/chi/v5"
)

type MedicalController struct {
	MedicalService interfaces.MedicalService
}

func NewMedicalController(medicalService interfaces.MedicalService) *MedicalController {
	return &MedicalController{MedicalService: medicalService}
}

func (c *MedicalController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)
	r.Post("/patient", c.handleAddMedicalPatient)

	return r
}

func (c *MedicalController) handleAddMedicalPatient(w http.ResponseWriter, r *http.Request) {
	payload := &medical_entity.AddMedicalPatient{}

	err := helpers.DecodeJSON(r, payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   err.Error(),
			Message: "Failed to decode JSON",
		})
		return
	}

	err = helpers.ValidatePayload(payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   err.Error(),
			Message: "Request doesn’t pass validation",
		})
		return
	}

	err = c.MedicalService.CreatePatient(r.Context(), payload)
	if errors.Is(err, medical_error.ErrIdentityNumberAlreadyExists) {
		helpers.ResponseJSON(w, http.StatusConflict, &helpers.ResponseBody{
			Error:   "Conflict error",
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, &helpers.ResponseBody{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	helpers.ResponseJSON(w, http.StatusCreated, &helpers.ResponseBody{
		Message: "Medical patient successfully added",
	})
}
