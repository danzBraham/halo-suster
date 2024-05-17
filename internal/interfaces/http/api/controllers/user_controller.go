package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	user_error "github.com/danzBraham/halo-suster/internal/exceptions/users"
	"github.com/danzBraham/halo-suster/internal/helpers"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	Service interfaces.UserService
	Router  chi.Router
}

func NewUserController(userService interfaces.UserService) *UserController {
	return &UserController{Service: userService}
}

func (c *UserController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/it/register", c.handleRegister)
	r.Post("/it/login", c.handleLogin)
	return r
}

func (c *UserController) handleRegister(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.RegisterITUser{}

	err := helpers.DecodeJSON(r, payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = helpers.ValidatePayload(payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Message: "Request doesn’t pass validation",
			Error:   err.Error(),
		})
		return
	}

	user, err := c.Service.CreateITUser(r.Context(), payload)
	if errors.Is(err, user_error.NIPAlreadyExistsError) {
		helpers.ResponseJSON(w, http.StatusConflict, &helpers.ResponseBody{
			Message: err.Error(),
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, &helpers.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   user.AccessToken,
		Expires: time.Now().Add(2 * time.Hour),
	}
	http.SetCookie(w, cookie)

	helpers.ResponseJSON(w, http.StatusCreated, &helpers.ResponseBody{
		Message: "User successfully registered",
		Data:    user,
	})
}

func (c *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.LoginUser{}

	err := helpers.DecodeJSON(r, payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = helpers.ValidatePayload(payload)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Message: "Request doesn’t pass validation",
			Error:   err.Error(),
		})
		return
	}

	user, err := c.Service.UserLogin(r.Context(), payload)
	if errors.Is(err, user_error.UserNotFoundError) {
		helpers.ResponseJSON(w, http.StatusConflict, &helpers.ResponseBody{
			Message: err.Error(),
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, &helpers.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   user.AccessToken,
		Expires: time.Now().Add(2 * time.Hour),
	}
	http.SetCookie(w, cookie)

	helpers.ResponseJSON(w, http.StatusCreated, &helpers.ResponseBody{
		Message: "User successfully login",
		Data:    user,
	})
}
