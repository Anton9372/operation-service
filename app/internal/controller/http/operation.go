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
	operationURL     = "/api/operations"
	operationByIdURL = "/api/operations/one/:uuid"
)

type OperationService interface {
	Create(ctx context.Context, dto dto.CreateOperationDTO) (string, error)
	GetByUUID(ctx context.Context, uuid string) (entity.Operation, error)
	Update(ctx context.Context, dto dto.UpdateOperationDTO) error
	Delete(ctx context.Context, uuid string) error
}

type operationHandler struct {
	service OperationService
	logger  *logging.Logger
}

func NewOperationHandler(service OperationService, logger *logging.Logger) Handler {
	return &operationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *operationHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, operationURL, apperror.Middleware(h.CreateOperation))
	router.HandlerFunc(http.MethodGet, operationByIdURL, apperror.Middleware(h.GetOperationByUUID))
	router.HandlerFunc(http.MethodPatch, operationByIdURL, apperror.Middleware(h.PartiallyUpdateOperation))
	router.HandlerFunc(http.MethodDelete, operationByIdURL, apperror.Middleware(h.DeleteOperation))
}

// CreateOperation
// @Summary 	Create operation
// @Description Creates new operation
// @Tags 		Operation
// @Accept		json
// @Param 		input	body 	 dto.CreateOperationDTO	true	"Operation's data"
// @Success 	201
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations [post]
func (h *operationHandler) CreateOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Create operation")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var createdOperation dto.CreateOperationDTO

	if err := json.NewDecoder(r.Body).Decode(&createdOperation); err != nil {
		h.logger.Error(err)
		return apperror.BadRequestError("invalid JSON body")
	}

	if createdOperation.CategoryUUID == "" || createdOperation.MoneySum == 0 {
		return apperror.BadRequestError("missing required fields")
	}

	operationUUID, err := h.service.Create(r.Context(), createdOperation)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", operationURL, operationUUID))
	w.WriteHeader(http.StatusCreated)

	h.logger.Info("Create operation successfully")
	return nil
}

// GetOperationByUUID
// @Summary 	Get operation by uuid
// @Description Get operation by uuid
// @Tags 		Operation
// @Produce 	json
// @Param 		uuid 	path 	 string 	true   "Operation's uuid"
// @Success 	200		{object} entity.Operation  "Operation"
// @Failure 	404 	{object} apperror.AppError "Operation not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router 		/operations/one/	[get]
func (h *operationHandler) GetOperationByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get operation by uuid")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")
	if operationUUID == "" {
		return apperror.BadRequestError("operation uuid must not be empty")
	}

	operation, err := h.service.GetByUUID(r.Context(), operationUUID)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, err = json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	h.logger.Info("Get operation by uuid successfully")
	return nil
}

// PartiallyUpdateOperation
// @Summary 	Update Operation
// @Description Update Operation
// @Tags 		Operation
// @Accept		json
// @Param 		uuid 		path 	 string 				true  "Operation's uuid"
// @Param 		input 		body 	 dto.UpdateOperationDTO true  "Operation's data"
// @Success 	204
// @Failure 	400 	{object} apperror.AppError "Validation error"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations/one [patch]
func (h *operationHandler) PartiallyUpdateOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Update operation")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")
	if operationUUID == "" {
		return apperror.BadRequestError("operation uuid must not be empty")
	}

	var updatedOperation dto.UpdateOperationDTO

	if err := json.NewDecoder(r.Body).Decode(&updatedOperation); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	updatedOperation.UUID = operationUUID

	err := h.service.Update(r.Context(), updatedOperation)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	h.logger.Info("Update operation successfully")
	return nil
}

// DeleteOperation
// @Summary 	Delete operation
// @Description Delete operation
// @Tags 		Operation
// @Param 		uuid 	path 	 string 	true  "Operation's uuid"
// @Success 	204
// @Failure 	404 	{object} apperror.AppError "Operation is not found"
// @Failure 	418 	{object} apperror.AppError "Something wrong with application logic"
// @Failure 	500 	{object} apperror.AppError "Internal server error"
// @Router /operations/one [delete]
func (h *operationHandler) DeleteOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Delete operation")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	operationUUID := params.ByName("uuid")
	if operationUUID == "" {
		return apperror.BadRequestError("operation uuid must not be empty")
	}

	err := h.service.Delete(r.Context(), operationUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	h.logger.Info("Delete operation successfully")
	return nil
}
