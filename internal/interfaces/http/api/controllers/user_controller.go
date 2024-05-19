package controllers

import (
	"errors"
	"net/http"
	"strconv"
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

	r.Post("/it/register", c.handleRegisterITUser)
	r.Post("/it/login", c.handleLoginITUser)
	r.Post("/nurse/login", c.handleLoginNurseUser)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/nurse/register", c.handleRegisterNurseUser)
		r.Get("/", c.handleGetUsers)
		r.Put("/nurse/{userId}", c.handleUpdateNurseUser)
		r.Delete("/nurse/{userId}", c.handleDeleteNurseUser)
		r.Post("/nurse/{userId}/access", c.handleGiveAccessNurseUser)
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

	if strconv.Itoa(payload.NIP)[:3] != "615" {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Validation error",
			Message: "Request doesn’t pass validation",
		})
		return
	}

	user, err := c.Service.CreateITUser(r.Context(), payload)
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
			Message: user_error.ErrUserIsNotIT.Error(),
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

	if strconv.Itoa(payload.NIP)[:3] != "303" {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Validation error",
			Message: "Request doesn’t pass validation",
		})
		return
	}

	user, err := c.Service.CreateNurseUser(r.Context(), payload)
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

func (c *UserController) handleLoginITUser(w http.ResponseWriter, r *http.Request) {
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

	if strconv.Itoa(payload.NIP)[:3] != "615" {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: user_error.ErrUserIsNotIT.Error(),
		})
		return
	}

	user, err := c.Service.LoginUser(r.Context(), payload)
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

func (c *UserController) handleLoginNurseUser(w http.ResponseWriter, r *http.Request) {
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

	if strconv.Itoa(payload.NIP)[:3] != "303" {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: user_error.ErrUserIsNotNurse.Error(),
		})
		return
	}

	user, err := c.Service.LoginUser(r.Context(), payload)
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

func (c *UserController) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &user_entity.UserQueryParams{
		UserID:    query.Get("userId"),
		Limit:     5,
		Offset:    0,
		NIP:       query.Get("nip"),
		Name:      query.Get("name"),
		Role:      query.Get("role"),
		CreatedAt: query.Get("createdAt"),
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = l
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = o
		}
	}

	users, err := c.Service.GetUsers(r.Context(), params)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusInternalServerError, &helpers.ResponseBody{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
		Message: "success",
		Data:    users,
	})
}

func (c *UserController) handleUpdateNurseUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	payload := &user_entity.UpdateNurseUser{
		UserID: userId,
	}

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

	if strconv.Itoa(payload.NIP)[:3] != "303" {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: user_error.ErrUserIsNotNurse.Error(),
		})
		return
	}

	err = c.Service.UpdateNurseUser(r.Context(), payload)
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

	helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
		Message: "Nurse user successfully updated",
	})
}

func (c *UserController) handleDeleteNurseUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	payload := &user_entity.UpdateNurseUser{
		UserID: userId,
	}

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

	if strconv.Itoa(payload.NIP)[:3] != "303" {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: user_error.ErrUserIsNotNurse.Error(),
		})
		return
	}

	err = c.Service.DeleteNurseUser(r.Context(), userId)
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

	helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
		Message: "Nurse user successfully deleted",
	})
}

func (c *UserController) handleGiveAccessNurseUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	payload := &user_entity.GiveAccessNurseUser{
		UserID: userId,
	}

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

	err = c.Service.GiveAccessNurseUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrUserNotFound) {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Error:   "Not found error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrUserIsNotNurse) {
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

	helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
		Message: "Nurse user successfully granted access",
	})
}
