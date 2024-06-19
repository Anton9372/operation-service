package service

import (
	"context"
	"fmt"
	"operation-service/internal/controller/dto"
	controller "operation-service/internal/controller/http"
	"operation-service/internal/domain/entity"
	"operation-service/pkg/logging"
)

type OperationRepo interface {
	Create(ctx context.Context, operation entity.Operation) (string, error)
	FindByUUID(ctx context.Context, uuid string) (entity.Operation, error)
	Update(ctx context.Context, operation entity.Operation) error
	Delete(ctx context.Context, uuid string) error
}

type operationService struct {
	repository OperationRepo
	logger     *logging.Logger
}

func NewOperationService(repository OperationRepo, logger *logging.Logger) controller.OperationService {
	return &operationService{
		repository: repository,
		logger:     logger,
	}
}

func (s *operationService) Create(ctx context.Context, dto dto.CreateOperationDTO) (string, error) {
	//todo validation

	operation := entity.NewOperation(dto)
	operationUUID, err := s.repository.Create(ctx, *operation)
	if err != nil {
		return "", fmt.Errorf("failed to create operation: %w", err)
	}
	return operationUUID, nil
}

func (s *operationService) GetByUUID(ctx context.Context, uuid string) (entity.Operation, error) {
	operation, err := s.repository.FindByUUID(ctx, uuid)
	if err != nil {
		return operation, fmt.Errorf("failed to find operation by uuid: %w", err)
	}
	return operation, nil
}

func (s *operationService) Update(ctx context.Context, dto dto.UpdateOperationDTO) error {
	//todo validation

	operation, err := s.repository.FindByUUID(ctx, dto.UUID)
	if err != nil {
		return fmt.Errorf("failed to find operation by uuid: %w", err)
	}

	updOperation := entity.UpdatedOperation(operation, dto)

	err = s.repository.Update(ctx, *updOperation)
	if err != nil {
		return fmt.Errorf("failed to update operation: %w", err)
	}
	return nil
}

func (s *operationService) Delete(ctx context.Context, uuid string) error {
	err := s.repository.Delete(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete operation by uuid: %w", err)
	}
	return nil
}
