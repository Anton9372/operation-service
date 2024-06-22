package dto

type CreateOperationDTO struct {
	CategoryUUID string  `json:"category_uuid"`
	MoneySum     float64 `json:"money_sum"`
	Description  string  `json:"description"`
}

type UpdateOperationDTO struct {
	UUID         string  `json:"uuid"`
	CategoryUUID string  `json:"category_uuid"`
	MoneySum     float64 `json:"money_sum"`
	Description  string  `json:"description"`
}
