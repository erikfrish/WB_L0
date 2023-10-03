package db

import (
	"L0/internal/config"
	// "L0/internal/storage"
	order_data "L0/internal/strct"
	"database/sql"

	// "errors"
	"fmt"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(storage config.Storage) (*Storage, error) {
	const op = "storage.db.New"
	driver_name := "pgx"
	storage_str := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		storage.Host, storage.Port, storage.User, storage.Password, storage.DB_Name)
	db, err := sql.Open(driver_name, storage_str)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "L0_orders" (
			id            VARCHAR(25) NOT NULL PRIMARY KEY,
			data          JSON
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(id string, data order_data.Data) (int64, error) {
	const op = "storage.db.Save"

	stmt, err := s.db.Prepare("INSERT INTO \"L0_orders\"(id, data) VALUES('$1','$2');")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement %w, %s", op, err, data)
	}
	data_value, _ := data.Value()
	res, err := stmt.Exec(id, data_value)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			return 0, fmt.Errorf("%s: %w", op, pgErr)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	idx, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last inserted id: %w", op, err)
	}
	return idx, nil
}

// func (s *Storage) GetURL(alias string) (string, error) {
// 	const op = "storage.db.GetURL"

// 	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = '?'")
// 	if err != nil {
// 		return "", fmt.Errorf("%s: prepare statement %w", op, err)
// 	}
// 	var resURL string
// 	err = stmt.QueryRow(alias).Scan((&resURL))
// 	if errors.Is(err, sql.ErrNoRows) {
// 		return "", storage.ErrURLNotFound
// 	}
// 	if err != nil {
// 		return "", fmt.Errorf("%s: execute statement: %w", op, err)
// 	}

// 	return resURL, nil
// }

// TODO: func (s *Storage) DeleteURL(alias string) error {}
