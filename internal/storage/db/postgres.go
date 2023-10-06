package db

import (
	// "context"
	order "L0/internal/strct"
	"L0/pkg/client/postgresql"

	// "L0/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	// "L0/pkg/logger/sl"
)

type Repository interface {
	Get(ctx context.Context, order_uid string) (any, error)
	GetAll(ctx context.Context) (any, error)
	Insert(ctx context.Context, data order.Data) error
	// And maybe more in future
}

type repository struct {
	client postgresql.Client
	logger *slog.Logger
}

func (r *repository) Get(ctx context.Context, order_uid string) (any, error) {
	q := `SELECT data FROM "L0_orders" WHERE id = $1`
	var row []byte
	err := r.client.QueryRow(ctx, q, order_uid).Scan(&row)
	if err != nil {
		return order.Data{}, fmt.Errorf("can't get data %e", err)
	}
	var res order.Data
	res.Scan(row)
	return res, nil
}

func (r *repository) GetAll(ctx context.Context) (any, error) {
	q := `SELECT data FROM "L0_orders"`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	res := make([]order.Data, 0)
	for rows.Next() {
		var row []byte

		err = rows.Scan(&row)
		if err != nil {
			return nil, err
		}
		var data_row order.Data
		data_row.Scan(row)
		res = append(res, data_row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) Insert(ctx context.Context, data order.Data) error {
	q := `INSERT INTO "L0_orders" (id, data) VALUES ($1, $2) RETURNING id`
	order_uid := data.OrderUID

	data_row, err := data.Value()
	if err != nil {
		return err
	}
	if _, err := r.client.Exec(ctx, q, order_uid, data_row); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return newErr
		}
		return errors.Join(errors.New("query error "), err)
	}
	return nil
}

func NewRepository(client postgresql.Client, logger *slog.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
