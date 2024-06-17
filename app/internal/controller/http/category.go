package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"opertion-service/internal/apperror"
	"opertion-service/internal/controller/dto"
	"opertion-service/pkg/logging"
)

const (
	categoryURL     = "/api/category"
	categoryBuIdURL = "/api/category/uuid/:uuid"
)

type CategoryService interface {
	Create(ctx context.Context, dto dto.CreateCategoryDTO) (string, error)
	Update(ctx context.Context, dto dto.UpdateCategoryDTO) error
	Delete(ctx context.Context, uuid string) error
}

type categoryHandler struct {
	service CategoryService
	logger  *logging.Logger
}

func NewCategoryHandler(service CategoryService, logger *logging.Logger) Handler {
	return &categoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, categoryURL, apperror.Middleware(h.CreateCategory))
	router.HandlerFunc(http.MethodPatch, categoryBuIdURL, apperror.Middleware(h.PartiallyUpdateCategory))
	router.HandlerFunc(http.MethodDelete, categoryBuIdURL, apperror.Middleware(h.DeleteCategory))
}

func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Create category")
	w.Header().Set("Content-Type", "application/json")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Fatalf("Error closing body %v", err)
		}
	}(r.Body)

	var createdCategory dto.CreateCategoryDTO

	if err := json.NewDecoder(r.Body).Decode(&createdCategory); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}
	//todo categoryType
	if createdCategory.Name == "" {
		return apperror.BadRequestError("missing required fields")
	}

	categoryUUID, err := h.service.Create(r.Context(), createdCategory)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", categoryURL, categoryUUID))
	w.WriteHeader(http.StatusCreated)

	h.logger.Info("Create category successfully")
	return nil
}

func (h *categoryHandler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Partially update category")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")

	defer func(Body io.ReadCloser) {
		err := r.Body.Close()
		if err != nil {
			h.logger.Fatalf("Error closing body %v", err)
		}
	}(r.Body)

	var updatedCategory dto.UpdateCategoryDTO

	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	updatedCategory.UUID = categoryUUID

	err := h.service.Update(r.Context(), updatedCategory)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	h.logger.Info("Update category successfully")
	return nil
}

func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Delete category")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")

	err := h.service.Delete(r.Context(), categoryUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	h.logger.Info("Delete category successfully")
	return nil
}
