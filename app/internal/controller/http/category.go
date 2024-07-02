package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"operation-service/internal/apperror"
	"operation-service/internal/controller/dto"
	"operation-service/internal/domain/entity"
	"operation-service/pkg/logging"
	"operation-service/pkg/utils"
)

const (
	categoryURL         = "/api/categories"
	categoryByIdURL     = "/api/categories/one/:uuid"
	categoryByUserIdURL = "/api/categories/user_uuid/:user_uuid"
)

type CategoryService interface {
	Create(ctx context.Context, dto dto.CreateCategoryDTO) (string, error)
	GetByUUID(ctx context.Context, uuid string) (entity.Category, error)
	GetByUserUUID(ctx context.Context, uuid string) ([]entity.Category, error)
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
	router.HandlerFunc(http.MethodGet, categoryByIdURL, apperror.Middleware(h.GetCategoryByUUID))
	router.HandlerFunc(http.MethodGet, categoryByUserIdURL, apperror.Middleware(h.GetCategoriesByUserUUID))
	router.HandlerFunc(http.MethodPatch, categoryByIdURL, apperror.Middleware(h.PartiallyUpdateCategory))
	router.HandlerFunc(http.MethodDelete, categoryByIdURL, apperror.Middleware(h.DeleteCategory))
}

// CreateCategory
// @Summary 	Create category
// @Description Creates new category
// @Tags 		Category
// @Accept		json
// @Param 		input	body 	 dto.CreateCategoryDTO	true	"Category data"
// @Success 	201
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories [post]
func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Create category")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var createdCategory dto.CreateCategoryDTO

	if err := json.NewDecoder(r.Body).Decode(&createdCategory); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	if createdCategory.UserUUID == "" || createdCategory.Name == "" || createdCategory.Type == "" {
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

// GetCategoryByUUID
// @Summary 	Get category by uuid
// @Description Get category by uuid
// @Tags 		Category
// @Produce 	json
// @Param 		uuid 	path 	 string 	true   "Category's uuid"
// @Success 	200		{object} entity.Category "Category"
// @Failure 	404 	{object} apperror.AppError "Category not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router 		/categories/one/	[get]
func (h *categoryHandler) GetCategoryByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get category by uuid")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")
	if categoryUUID == "" {
		return apperror.BadRequestError("category uuid must not be empty")
	}

	category, err := h.service.GetByUUID(r.Context(), categoryUUID)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, err = json.Marshal(category)
	if err != nil {
		return fmt.Errorf("failed to marshal category: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	h.logger.Info("Get category by uuid successfully")
	return nil
}

// GetCategoriesByUserUUID
// @Summary 	Get categories by user's uuid
// @Description Get list of categories belonging to user
// @Tags 		Category
// @Produce 	json
// @Param 		user_uuid 	path 	 string 	true   "User's uuid"
// @Success 	200			{object} []entity.Category "Categories"
// @Failure 	404 		{object} apperror.AppError "User not found"
// @Failure 	418 		{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 		{object} apperror.AppError "Internal server error"
// @Router 		/categories/user_uuid/	[get]
func (h *categoryHandler) GetCategoriesByUserUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get categories by user's uuid")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("user_uuid")
	if userUUID == "" {
		return apperror.BadRequestError("user's uuid must not be empty")
	}

	categories, err := h.service.GetByUserUUID(r.Context(), userUUID)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, err = json.Marshal(categories)
	if err != nil {
		return fmt.Errorf("failed to marshal categories: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	h.logger.Info("Get categories by user's uuid successfully")
	return nil
}

// PartiallyUpdateCategory
// @Summary 	Update category
// @Description Update category
// @Tags 		Category
// @Accept		json
// @Param 		uuid 		path 	 string 				true  "Category's uuid"
// @Param 		input 		body 	 dto.UpdateCategoryDTO true  "Category's data"
// @Success 	204
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories/one [patch]
func (h *categoryHandler) PartiallyUpdateCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Partially update category")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")
	if categoryUUID == "" {
		return apperror.BadRequestError("category uuid must not be empty")
	}

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

// DeleteCategory
// @Summary 	Delete category
// @Description Delete category
// @Tags 		Category
// @Param 		uuid 	path 	 string 	true  "Category's uuid"
// @Success 	204
// @Failure 	404 	{object} apperror.AppError "Category is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /categories/one [delete]
func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Delete category")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	categoryUUID := params.ByName("uuid")
	if categoryUUID == "" {
		return apperror.BadRequestError("category uuid must not be empty")
	}

	err := h.service.Delete(r.Context(), categoryUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	h.logger.Info("Delete category successfully")
	return nil
}
