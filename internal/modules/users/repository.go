package users

import (
	"context"
	"encoding/json"
	"time"

	"github.com/username/project-name/ent"
	"github.com/username/project-name/internal/platform/cache"
)

type Repository struct {
	db    *ent.Client
	cache cache.Cache
}

func NewRepository(db *ent.Client, cache cache.Cache) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]*ent.User, error) {
	return r.db.User.Query().AllX(ctx), nil
}

func (r *Repository) GetUser(
	ctx context.Context,
	userId string,
) (*ent.User, error) {
	cUser, err := r.cache.Get(ctx, userId)

	if err == nil && cUser != "" {
		var user ent.User
		if err := json.Unmarshal([]byte(cUser), &user); err == nil {
			return &user, nil
		}
	}

	user, err := r.db.User.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	r.cache.Set(
		ctx,
		user.ID,
		data,
		int(15*time.Minute.Seconds()),
	)

	return user, nil

}
