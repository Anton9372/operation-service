package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"operation-service/internal/apperror"
	"operation-service/internal/domain/entity"
	"operation-service/internal/domain/service"
	"operation-service/pkg/logging"
	"operation-service/pkg/postgresql"
	"operation-service/pkg/utils"
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

func handleSQLError(err error, logger *logging.Logger) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
		logger.Error(newErr)

		if pgErr.Code == "23505" { //uniqueness violation
			return apperror.BadRequestError("This category name already exists")
		}
		return newErr
	}

	return err
}

func (r *categoryRepo) Create(ctx context.Context, category entity.Category) (string, error) {
	query := `
				INSERT INTO categories
					(name, type)
				VALUES
					($1, $2)
				RETURNING id;
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	var categoryUUID string
	err := r.client.QueryRow(nCtx, query, category.Name, category.Type).Scan(&categoryUUID)
	if err != nil {
		return "", handleSQLError(err, r.logger)
	}

	return categoryUUID, nil
}

func (r *categoryRepo) FindAll(ctx context.Context) ([]entity.Category, error) {
	query := `
				SELECT
					id, name, type
				FROM
					categories
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	rows, err := r.client.Query(nCtx, query)
	if err != nil {
		return nil, handleSQLError(err, r.logger)
	}
	defer rows.Close()

	categories := make([]entity.Category, 0)
	for rows.Next() {
		var category entity.Category
		err = rows.Scan(&category.UUID, &category.Name, &category.Type)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) FindByUUID(ctx context.Context, uuid string) (entity.Category, error) {
	query := `
				SELECT
					id, name, type
				FROM
					categories
				WHERE
					id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	var category entity.Category
	err := r.client.QueryRow(nCtx, query, uuid).Scan(&category.UUID, &category.Name, &category.Type)
	if err != nil {
		return entity.Category{}, handleSQLError(err, r.logger)
	}

	return category, nil
}

func (r *categoryRepo) FindByName(ctx context.Context, name string) (entity.Category, error) {
	query := `
				SELECT
					id, name, type
				FROM
		    		categories
				WHERE
		    		name = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", utils.FormatSQLQuery(query)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	var category entity.Category
	err := r.client.QueryRow(nCtx, query, name).Scan(&category.UUID, &category.Name, &category.Type)
	if err != nil {
		return entity.Category{}, handleSQLError(err, r.logger)
	}

	return category, nil
}

func (r *categoryRepo) Update(ctx context.Context, category entity.Category) error {
	query := `
				UPDATE 
					categories
				SET 
    				name = $1
				WHERE
				    id = $2
	`
	r.logger.Trace("SQL query: %s", utils.FormatSQLQuery(query))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()

	cmdTag, err := r.client.Exec(nCtx, query, category.Name, category.UUID)
	if err != nil {
		return handleSQLError(err, r.logger)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, uuid string) error {
	//TODO
	query := `
				DELETE FROM
			    	categories
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
