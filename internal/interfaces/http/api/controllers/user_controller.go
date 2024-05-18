package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/danzBraham/halo-suster/internal/applications/interfaces"
	user_entity "github.com/danzBraham/halo-suster/internal/domains/entities/users"
	user_error "github.com/danzBraham/halo-suster/internal/exceptions/users"
	"github.com/danzBraham/halo-suster/internal/helpers"
	"github.com/danzBraham/halo-suster/internal/interfaces/http/api/middlewares"
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

	r.Route("/it", func(r chi.Router) {
		r.Post("/register", c.handleRegisterITUser)
		r.Post("/login", c.handleLogin)
	})

	r.Route("/nurse", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/register", c.handleRegisterNurseUser)
	})

	return r
}

func (c *UserController) handleRegisterITUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.RegisterITUser{}

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
			Error:   "Validation error",
			Message: "Request doesn’t pass validation",
		})
		return
	}

	user, err := c.Service.CreateITUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrNotITUserNIP) {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrUserNotFound) {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrNIPAlreadyExists) {
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

func (c *UserController) handleRegisterNurseUser(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middlewares.ContextRoleKey)
	if role != user_entity.IT {
		helpers.ResponseJSON(w, http.StatusUnauthorized, &helpers.ResponseBody{
			Error:   "Unauthorized error",
			Message: "You are not IT user",
		})
		return
	}

	payload := &user_entity.RegisterNurseUser{}

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
			Error:   "Validation error",
			Message: "Request doesn’t pass validation",
		})
		return
	}

	user, err := c.Service.CreateNurseUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrNotNurseUserNIP) {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrUserNotFound) {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrNIPAlreadyExists) {
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

	user, err := c.Service.UserLogin(r.Context(), payload)
	if errors.Is(err, user_error.ErrInvalidPassword) {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Bad request error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrUserNotFound) {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
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
