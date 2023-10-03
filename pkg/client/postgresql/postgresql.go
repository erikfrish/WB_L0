package postgresql

import (
	"context"
	"fmt"
	"log"

	"L0/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"L0/pkg/logger/sl"
	repeatable "L0/pkg/utils"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// TODO: add slog logging
func NewClient(ctx context.Context, maxAttempts int, sc config.Storage) (pool *pgxpool.Pool, err error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		sc.Host, sc.Port, sc.User, sc.Password, sc.DB_Name)
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, connString)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql", sl.Err(err))
	}

	return pool, nil
}
