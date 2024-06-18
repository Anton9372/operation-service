package entity

import (
	"operation-service/internal/controller/dto"
	"time"
)

type Operation struct {
	UUID         string    `json:"uuid"`
	CategoryUUID string    `json:"category_uuid"`
	Description  string    `json:"description"`
	MoneySum     float64   `json:"money_sum"`
	DateTime     time.Time `json:"date_time"`
}

func NewOperation(dto dto.CreateOperationDTO) *Operation {
	return &Operation{
		CategoryUUID: dto.CategoryUUID,
		Description:  dto.Description,
		MoneySum:     dto.MoneySum,
		DateTime:     dto.DateTime,
	}
}

func UpdatedOperation(existing Operation, dto dto.UpdateOperationDTO) *Operation {
	updOperation := &Operation{
		UUID:     dto.UUID,
		DateTime: existing.DateTime,
	}
	if dto.Description != "" {
		updOperation.Description = dto.Description
	}
}
