package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB interface {
	Select(ctx context.Context, ptr any, query string, args ...any) error
	SelectRow(ctx context.Context, ptr any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	ExecNamed(ctx context.Context, query string, arg any) (sql.Result, error)
}

type Client struct {
	*sqlx.DB
}

func Open(ctx context.Context, dsn string) (*Client, error) {
	db := &Client{}

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "parse config")
	}

	sqlDb := stdlib.OpenDB(*cfg)

	pgDb := sqlx.NewDb(sqlDb, "pgx")
	err = pgDb.PingContext(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "ping database")
	}

	db.DB = pgDb
	db.DB.SetMaxOpenConns(100)

	return db, nil
}

func (db *Client) CreateSchema(schema string) error {
	_, err := db.DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		return errors.WithMessage(err, "exec query")
	}
	return nil
}

func (db *Client) Select(ctx context.Context, ptr any, query string, args ...any) error {
	return db.SelectContext(ctx, ptr, query, args...)
}

func (db *Client) SelectRow(ctx context.Context, ptr any, query string, args ...any) error {
	return db.GetContext(ctx, ptr, query, args...)
}

func (db *Client) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.ExecContext(ctx, query, args...)
}

func (db *Client) ExecNamed(ctx context.Context, query string, arg any) (sql.Result, error) {
	return db.NamedExecContext(ctx, query, arg)
}
