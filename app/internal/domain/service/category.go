package service

import (
	"context"
	"opertion-service/internal/controller/dto"
	"opertion-service/internal/controller/http"
	"opertion-service/internal/domain/entity"
	"opertion-service/pkg/logging"
)

type CategoryRepo interface {
	Create(ctx context.Context, category entity.Category) (string, error)
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, uuid string) error
}

type categoryService struct {
	repository CategoryRepo
	logger     *logging.Logger
}

func NewCategoryService(repository CategoryRepo, logger *logging.Logger) http.CategoryService {
	return &categoryService{
		repository: repository,
		logger:     logger,
	}
}

func (c categoryService) Create(ctx context.Context, dto dto.CreateCategoryDTO) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c categoryService) Update(ctx context.Context, dto dto.UpdateCategoryDTO) error {
	//TODO implement me
	panic("implement me")
}

func (c categoryService) Delete(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}
