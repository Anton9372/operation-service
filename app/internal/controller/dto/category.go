package dto

type CreateCategoryDTO struct {
	Name string `json:"name"`
	//Type string `json:"type"`
}

type UpdateCategoryDTO struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
