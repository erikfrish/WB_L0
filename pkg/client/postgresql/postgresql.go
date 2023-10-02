package postgresql

import (
	"context"
	"fmt"
	"log"

	"L0/internal/config"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repeatable "L0/pkg/utils"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.Storage) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", sc.User, sc.Password, sc.Host, sc.Port, sc.DB_Name)
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}
