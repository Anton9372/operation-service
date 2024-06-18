package dto

import "time"

type CreateOperationDTO struct {
	CategoryUUID string    `json:"category_uuid"`
	Description  string    `json:"description"`
	MoneySum     float64   `json:"money_sum"`
	DateTime     time.Time `json:"date_time"`
}

type UpdateOperationDTO struct {
	UUID         string  `json:"uuid"`
	CategoryUUID string  `json:"category_uuid"`
	Description  string  `json:"description"`
	MoneySum     float64 `json:"money_sum"`
}
