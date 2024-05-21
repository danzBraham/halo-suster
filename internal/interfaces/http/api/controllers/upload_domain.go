package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	upload_entity "github.com/danzBraham/halo-suster/internal/domains/entities/uploads"
	"github.com/danzBraham/halo-suster/internal/helpers"
	"github.com/danzBraham/halo-suster/internal/interfaces/http/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)
	r.Post("/image", c.handleUploadImage)

	return r
}

func (c *UploadController) handleUploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(upload_entity.MaxUploadSize)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   err.Error(),
			Message: "Unable to parse form",
		})
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   err.Error(),
			Message: "Unable to get file from form",
		})
		return
	}
	defer file.Close()

	// Validate file type
	fileType := strings.ToLower(filepath.Ext(header.Filename))
	if fileType != ".jpg" && fileType != ".jpeg" {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Bad request error",
			Message: "Invalid file type. Only .jpg and .jpeg are allowed",
		})
		return
	}

	// Check file size
	if header.Size < upload_entity.MinUploadSize || header.Size > upload_entity.MaxUploadSize {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Bad request error",
			Message: "Invalid file size. Must be between 10KB and 2MB",
		})
		return
	}

	// Generate UUID for the new file name
	newFilename := uuid.New().String() + fileType
	filePath := filepath.Join(upload_entity.UploadPath, newFilename)

	// Create upload directory if not exists
	if _, err := os.Stat(upload_entity.UploadPath); os.IsNotExist(err) {
		err := os.MkdirAll(upload_entity.UploadPath, os.ModePerm)
		if err != nil {
			helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
				Error:   "Bad request error",
				Message: "Unable to create uplaod directory",
			})
			return
		}
	}

	// Save th file
	dst, err := os.Create(filePath)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Bad request error",
			Message: "Unable to save file",
		})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		helpers.ResponseJSON(w, http.StatusBadRequest, &helpers.ResponseBody{
			Error:   "Bad request error",
			Message: "Unable to save file",
		})
		return
	}

	helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
		Message: "image uploaded successfully",
		Data: &upload_entity.UploadedImage{
			ImageURL: fmt.Sprintf("http://localhost:8080/uploads/%s", newFilename),
		},
	})
}
