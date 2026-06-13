package auth

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/username/project-name/ent"
	"github.com/username/project-name/ent/refreshsession"
	"github.com/username/project-name/ent/user"
)

type Repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, role string, permissions []string) (*UserRecord, error) {
	created, err := r.db.User.Create().
		SetEmail(email).
		SetPasswordHash(passwordHash).
		SetRole(role).
		SetPermissions(permissions).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(created), nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	found, err := r.db.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(found), nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*UserRecord, error) {
	found, err := r.db.User.Query().Where(user.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(found), nil
}

func (r *Repository) GetUserByResetTokenHash(ctx context.Context, tokenHash string) (*UserRecord, error) {
	found, err := r.db.User.Query().Where(user.ResetTokenHashEQ(tokenHash)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(found), nil
}

func (r *Repository) SetResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	return r.db.User.UpdateOneID(userID).
		SetResetTokenHash(tokenHash).
		SetResetTokenExpiresAt(expiresAt).
		Exec(ctx)
}

func (r *Repository) ResetPassword(ctx context.Context, userID, passwordHash string) error {
	return r.db.User.UpdateOneID(userID).
		SetPasswordHash(passwordHash).
		ClearResetTokenHash().
		ClearResetTokenExpiresAt().
		Exec(ctx)
}

func (r *Repository) CreateRefreshSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (*RefreshSessionRecord, error) {
	session, err := r.db.RefreshSession.Create().
		SetUserID(userID).
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapSession(session), nil
}

func (r *Repository) GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (*RefreshSessionRecord, error) {
	session, err := r.db.RefreshSession.Query().Where(refreshsession.TokenHashEQ(tokenHash)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapSession(session), nil
}

func (r *Repository) RevokeRefreshSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	return r.db.RefreshSession.UpdateOneID(sessionID).SetRevokedAt(revokedAt).Exec(ctx)
}

func (r *Repository) RevokeUserRefreshSessions(ctx context.Context, userID string, revokedAt time.Time) error {
	_, err := r.db.RefreshSession.Update().
		Where(
			refreshsession.UserIDEQ(userID),
			refreshsession.RevokedAtIsNil(),
		).
		SetRevokedAt(revokedAt).
		Save(ctx)
	return err
}

func (r *Repository) IsEmailTakenError(err error) bool {
	return ent.IsConstraintError(err) || sqlgraph.IsConstraintError(err)
}

func mapUser(entity *ent.User) *UserRecord {
	if entity == nil {
		return nil
	}

	return &UserRecord{
		ID:                entity.ID,
		Email:             entity.Email,
		PasswordHash:      entity.PasswordHash,
		Role:              entity.Role,
		Permissions:       entity.Permissions,
		ResetTokenHash:    entity.ResetTokenHash,
		ResetTokenExpires: entity.ResetTokenExpiresAt,
	}
}

func mapSession(entity *ent.RefreshSession) *RefreshSessionRecord {
	if entity == nil {
		return nil
	}

	return &RefreshSessionRecord{
		ID:        entity.ID,
		UserID:    entity.UserID,
		TokenHash: entity.TokenHash,
		ExpiresAt: entity.ExpiresAt,
		RevokedAt: entity.RevokedAt,
	}
}
