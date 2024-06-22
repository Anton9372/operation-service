package service

import (
	"context"
	"fmt"
	"operation-service/internal/apperror"
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
	operationRepo OperationRepo
	categoryRepo  CategoryRepo
	logger        *logging.Logger
}

func NewOperationService(operationRepo OperationRepo, categoryRepo CategoryRepo,
	logger *logging.Logger) controller.OperationService {
	return &operationService{
		operationRepo: operationRepo,
		categoryRepo:  categoryRepo,
		logger:        logger,
	}
}

func (s *operationService) Create(ctx context.Context, dto dto.CreateOperationDTO) (string, error) {
	if dto.MoneySum <= 0 {
		return "", apperror.BadRequestError("money sum can not be negative or zero")
	}

	_, err := s.categoryRepo.FindByUUID(ctx, dto.CategoryUUID)
	if err != nil {
		return "", err
	}

	operation := entity.NewOperation(dto)

	operationUUID, err := s.operationRepo.Create(ctx, *operation)
	if err != nil {
		return "", fmt.Errorf("failed to create operation: %w", err)
	}
	return operationUUID, nil
}

func (s *operationService) GetByUUID(ctx context.Context, uuid string) (entity.Operation, error) {
	operation, err := s.operationRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return operation, fmt.Errorf("failed to find operation by uuid: %w", err)
	}
	return operation, nil
}

func (s *operationService) Update(ctx context.Context, dto dto.UpdateOperationDTO) error {
	if dto.MoneySum < 0 {
		return apperror.BadRequestError("money sum can not be negative")
	}

	operation, err := s.operationRepo.FindByUUID(ctx, dto.UUID)
	if err != nil {
		return fmt.Errorf("failed to find operation by uuid: %w", err)
	}

	updOperation := entity.UpdatedOperation(operation, dto)

	_, err = s.categoryRepo.FindByUUID(ctx, updOperation.CategoryUUID)
	if err != nil {
		return err
	}

	err = s.operationRepo.Update(ctx, *updOperation)
	if err != nil {
		return fmt.Errorf("failed to update operation: %w", err)
	}
	return nil
}

func (s *operationService) Delete(ctx context.Context, uuid string) error {
	_, err := s.operationRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.operationRepo.Delete(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete operation by uuid: %w", err)
	}
	return nil
}
