package postgres

import (
	"context"
	"fmt"
	"opertion-service/internal/domain/entity"
	"opertion-service/internal/domain/service"
	"opertion-service/pkg/logging"
	"opertion-service/pkg/postgresql"
	"opertion-service/pkg/utils"
	"time"
)

const queryWaitTime = 5 * time.Second

type categoryRepo struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewCategoryRepo(client postgresql.Client, logger *logging.Logger) service.CategoryRepo {
	return &categoryRepo{
		client: client,
		logger: logger,
	}
}

func (r *categoryRepo) Create(ctx context.Context, category entity.Category) (string, error) {
	query := ``
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	var categoryUUID string
	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()
	err := r.client.QueryRow(nCtx, query, category).Scan(&categoryUUID)
	if err != nil {
		return "", utils.HandleSQLError(err, r.logger)
	}

	return categoryUUID, nil
}

func (r *categoryRepo) Update(ctx context.Context, category entity.Category) error {
	//TODO implement me
	panic("implement me")
}

func (r *categoryRepo) Delete(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}
