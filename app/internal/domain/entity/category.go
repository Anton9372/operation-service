package entity

import (
	"operation-service/internal/controller/dto"
	"operation-service/internal/domain/types"
)

type Category struct {
	UUID string             `json:"uuid"`
	Name string             `json:"name"`
	Type types.CategoryType `json:"type"`
}

func NewCategory(dto dto.CreateCategoryDTO) *Category {
	return &Category{
		Name: dto.Name,
		Type: dto.Type,
	}
}

func UpdatedCategory(existing Category, dto dto.UpdateCategoryDTO) *Category {
	updCategory := &Category{
		UUID: dto.UUID,
		Type: existing.Type,
	}

	if dto.Name != "" {
		updCategory.Name = dto.Name
	} else {
		updCategory.Name = existing.Name
	}

	return updCategory
}
