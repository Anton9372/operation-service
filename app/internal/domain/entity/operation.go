package entity

import (
	"operation-service/internal/controller/dto"
	"time"
)

type Operation struct {
	UUID         string    `json:"uuid"`
	CategoryUUID string    `json:"category_uuid"`
	MoneySum     float64   `json:"money_sum"`
	Description  string    `json:"description"`
	DateTime     time.Time `json:"date_time"`
}

func NewOperation(dto dto.CreateOperationDTO) *Operation {
	return &Operation{
		CategoryUUID: dto.CategoryUUID,
		MoneySum:     dto.MoneySum,
		Description:  dto.Description,
		DateTime:     time.Now(),
	}
}

func UpdatedOperation(existing Operation, dto dto.UpdateOperationDTO) *Operation {
	updOperation := new(Operation)

	updOperation.UUID = dto.UUID

	if dto.CategoryUUID != "" {
		updOperation.CategoryUUID = dto.CategoryUUID
	} else {
		updOperation.CategoryUUID = existing.CategoryUUID
	}

	if dto.MoneySum != 0 {
		updOperation.MoneySum = dto.MoneySum
	} else {
		updOperation.MoneySum = existing.MoneySum
	}

	if dto.Description != "" {
		updOperation.Description = dto.Description
	} else {
		updOperation.Description = existing.Description
	}

	updOperation.DateTime = existing.DateTime

	return updOperation
}
