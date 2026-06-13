package db

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/lib/pq"
	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/config"
)

type Client struct {
	Ent *ent.Client
	sql *sql.DB
}

func New(ctx context.Context, cfg config.DatabaseConfig) (*Client, error) {
	if _, err := pq.ParseURL(cfg.URL); err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	sqlDB, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	driver := entsql.OpenDB(dialect.Postgres, sqlDB)

	return &Client{
		Ent: ent.NewClient(ent.Driver(driver)),
		sql: sqlDB,
	}, nil
}

func (c *Client) PingContext(ctx context.Context) error {
	return c.sql.PingContext(ctx)
}

func (c *Client) Close() error {
	return c.sql.Close()
}
