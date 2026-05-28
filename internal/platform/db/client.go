package db

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/config"

	_ "github.com/lib/pq"
)

func New(databaseUrl string, cfg *config.Config) (*ent.Client, error) {
	client, err := ent.Open(
		dialect.Postgres,
		databaseUrl,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"Failed opening database connection: %w",
			err,
		)
	}

	if cfg.AutoMigrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf(
				"Failed creating schema resources: %w",
				err,
			)
		}
	}

	return client, nil
}
