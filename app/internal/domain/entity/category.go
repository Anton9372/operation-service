package entity

import "opertion-service/internal/controller/dto"

//type CategoryType int
//
//const (
//	Income CategoryType = iota
//	Expence
//)

type Category struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	//Type CategoryType `json:"type"`
}

func NewCategory(dto dto.CreateCategoryDTO) *Category {
	return &Category{
		Name: dto.Name,
	}
}

func UpdatedCategory(dto dto.UpdateCategoryDTO) *Category {
	return &Category{
		UUID: dto.UUID,
		Name: dto.Name,
	}
}
