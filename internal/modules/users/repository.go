package users

import (
	"context"

	"github.com/username/project-name/ent"
)

type Repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]*ent.User, error) {
	return r.db.User.Query().AllX(ctx), nil
}
