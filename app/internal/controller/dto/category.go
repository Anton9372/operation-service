package dto

import "operation-service/internal/domain/types"

type CreateCategoryDTO struct {
	Name string             `json:"name"`
	Type types.CategoryType `json:"type"`
}

type UpdateCategoryDTO struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
