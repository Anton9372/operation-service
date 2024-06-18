package service

import (
	"context"
	"fmt"
	"operation-service/internal/apperror"
	"operation-service/internal/controller/dto"
	"operation-service/internal/controller/http"
	"operation-service/internal/domain/entity"
	"operation-service/internal/domain/types"
	"operation-service/pkg/logging"
)

type CategoryRepo interface {
	Create(ctx context.Context, category entity.Category) (string, error)
	FindAll(ctx context.Context) ([]entity.Category, error)
	FindByUUID(ctx context.Context, uuid string) (entity.Category, error)
	FindByName(ctx context.Context, name string) (entity.Category, error)
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, uuid string) error
}

type categoryService struct {
	repository CategoryRepo
	logger     *logging.Logger
}

func NewCategoryService(repository CategoryRepo, logger *logging.Logger) controller.CategoryService {
	return &categoryService{
		repository: repository,
		logger:     logger,
	}
}

func (s *categoryService) Create(ctx context.Context, dto dto.CreateCategoryDTO) (string, error) {
	if dto.Type != types.IncomeType && dto.Type != types.ExpenseType {
		return "", apperror.BadRequestError("category type must be 'Income' or 'Expense'")
	}

	category := entity.NewCategory(dto)
	categoryUUID, err := s.repository.Create(ctx, *category)
	if err != nil {
		return categoryUUID, fmt.Errorf("failed to create category: %w", err)
	}
	return categoryUUID, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repository.FindAll(ctx)
	if err != nil {
		return categories, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

func (s *categoryService) GetByUUID(ctx context.Context, uuid string) (entity.Category, error) {
	category, err := s.repository.FindByUUID(ctx, uuid)
	if err != nil {
		return category, fmt.Errorf("failed to get category by uuid: %w", err)
	}
	return category, nil
}

func (s *categoryService) GetByName(ctx context.Context, name string) (entity.Category, error) {
	category, err := s.repository.FindByName(ctx, name)
	if err != nil {
		return category, fmt.Errorf("failed to get category by name: %w", err)
	}
	return category, nil
}

func (s *categoryService) Update(ctx context.Context, dto dto.UpdateCategoryDTO) error {
	if dto.Name == "" {
		return apperror.BadRequestError("category name must not be empty")
	}

	category, err := s.repository.FindByUUID(ctx, dto.UUID)
	if err != nil {
		return fmt.Errorf("failed to find category bu uuid: %w", err)
	}

	updCategory := entity.UpdatedCategory(category, dto)

	err = s.repository.Update(ctx, *updCategory)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (s *categoryService) Delete(ctx context.Context, uuid string) error {
	err := s.repository.Delete(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return err
}
