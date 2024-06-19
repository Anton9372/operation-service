package postgres

import (
	"context"
	"fmt"
	"operation-service/internal/domain/entity"
	"operation-service/internal/domain/service"
	"operation-service/pkg/logging"
	"operation-service/pkg/postgresql"
	"operation-service/pkg/utils"
)

type operationRepo struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewOperationRepo(client postgresql.Client, logger *logging.Logger) service.OperationRepo {
	return &operationRepo{
		client: client,
		logger: logger,
	}
}

func (r *operationRepo) Create(ctx context.Context, operation entity.Operation) (string, error) {
	query := `
				INSERT INTO operations
					(category_id, money_sum, description, date_time)
				VALUES 
					($1, $2, $3, $4)
				RETURNING id;
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	var operationUUID string
	err := r.client.QueryRow(nCtx, query, operation.CategoryUUID, operation.MoneySum, operation.Description,
		operation.DateTime).Scan(&operationUUID)
	if err != nil {
		return "", handleSQLError(err, r.logger)
	}

	return operationUUID, nil
}

func (r *operationRepo) FindByUUID(ctx context.Context, uuid string) (entity.Operation, error) {
	query := `
				SELECT
					id, category_id, money_sum, description, date_time
				FROM
					operations
				WHERE
					id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	var operation entity.Operation
	err := r.client.QueryRow(nCtx, query, uuid).Scan(&operation.UUID, operation.CategoryUUID, operation.MoneySum,
		operation.Description, operation.DateTime)
	if err != nil {
		return entity.Operation{}, handleSQLError(err, r.logger)
	}

	return operation, nil
}

func (r *operationRepo) Update(ctx context.Context, operation entity.Operation) error {
	query := `
				UPDATE
					operations
				SET
					category_id = $1, money_sum = $2, description = $3
				WHERE
					id = $4 
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	cmdTag, err := r.client.Exec(nCtx, query, operation.CategoryUUID, operation.MoneySum, operation.Description,
		operation.UUID)
	if err != nil {
		return handleSQLError(err, r.logger)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil
}

func (r *operationRepo) Delete(ctx context.Context, uuid string) error {
	query := `
				DELETE FROM
					operations
				WHERE
				    id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	cmdTag, err := r.client.Exec(nCtx, query, uuid)
	if err != nil {
		return handleSQLError(err, r.logger)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows were deleted")
	}
	return nil
}
