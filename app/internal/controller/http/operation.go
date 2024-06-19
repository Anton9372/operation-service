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
	operationURL     = "/api/operation"
	operationByIdURL = "/api/operation/uuid/:uuid"
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
	//TODO implement me
	panic("implement me")
}

func (h *operationHandler) CreateOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Create operation")
	defer utils.CloseRequestBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var createdOperation dto.CreateOperationDTO

	if err := json.NewDecoder(r.Body).Decode(&createdOperation); err != nil {
		return apperror.BadRequestError("invalid JSON body")
	}

	//todo validation

	operationUUID, err := h.service.Create(r.Context(), createdOperation)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s:%s", operationURL, operationUUID))
	w.WriteHeader(http.StatusCreated)

	h.logger.Info("Create operation successfully")
	return nil
}

func (h *operationHandler) GetOperationByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get operation by uuid")
	defer utils.CloseRequestBody(h.logger, r.Body)
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

func (h *operationHandler) PartiallyUpdateOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Update operation")
	defer utils.CloseRequestBody(h.logger, r.Body)
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

func (h *operationHandler) DeleteOperation(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Delete operation")
	defer utils.CloseRequestBody(h.logger, r.Body)
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
