package users

import (
	"context"

	"github.com/username/project-name/ent"
)

type Repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]UserRecord, error) {
	users, err := r.db.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]UserRecord, 0, len(users))
	for _, user := range users {
		result = append(result, UserRecord{
			ID:    user.ID,
			Email: user.Email,
		})
	}

	return result, nil
}
