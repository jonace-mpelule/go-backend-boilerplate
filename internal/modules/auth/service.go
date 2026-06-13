package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/username/project-name/ent"
	apperrors "github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/platform/analytics"
	"github.com/username/project-name/internal/platform/mailer"
	"github.com/username/project-name/internal/utils"
)

type jwtGenerator interface {
	Generate(userID, role string, permissions []string) (string, error)
}

type passwordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) error
}

type mailerSender interface {
	Send(ctx context.Context, to, subject, html string) error
}

type authRepository interface {
	CreateUser(ctx context.Context, email, passwordHash, role string, permissions []string) (*UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*UserRecord, error)
	GetUserByID(ctx context.Context, id string) (*UserRecord, error)
	GetUserByResetTokenHash(ctx context.Context, tokenHash string) (*UserRecord, error)
	SetResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	ResetPassword(ctx context.Context, userID, passwordHash string) error
	CreateRefreshSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (*RefreshSessionRecord, error)
	GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (*RefreshSessionRecord, error)
	RevokeRefreshSession(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeUserRefreshSessions(ctx context.Context, userID string, revokedAt time.Time) error
	IsEmailTakenError(err error) bool
}

type Service struct {
	repo            authRepository
	jwt             jwtGenerator
	passwords       passwordHasher
	mailer          mailerSender
	analytics       analytics.Analytics
	tokenTTL        time.Duration
	refreshTokenTTL time.Duration
	resetTokenTTL   time.Duration
}

func NewService(
	repo authRepository,
	jwt jwtGenerator,
	passwords passwordHasher,
	mailerClient mailer.Mailer,
	analyticsClient analytics.Analytics,
	tokenTTL, refreshTokenTTL, resetTokenTTL time.Duration,
) *Service {
	return &Service{
		repo:            repo,
		jwt:             jwt,
		passwords:       passwords,
		mailer:          mailerClient,
		analytics:       analyticsClient,
		tokenTTL:        tokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		resetTokenTTL:   resetTokenTTL,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterRequest) (*AuthResponse, *apperrors.AppError) {
	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password")
	}

	created, err := s.repo.CreateUser(ctx, input.Email, passwordHash, "user", []string{})
	if err != nil {
		if s.repo.IsEmailTakenError(err) {
			return nil, apperrors.Conflict("email is already registered")
		}
		return nil, apperrors.Internal("failed to create user")
	}

	response, appErr := s.issueTokens(ctx, created)
	if appErr != nil {
		return nil, appErr
	}

	s.analytics.Track(ctx, "auth.registered", map[string]any{"user_id": created.ID})

	return response, nil
}

func (s *Service) Login(ctx context.Context, input LoginRequest) (*AuthResponse, *apperrors.AppError) {
	userRecord, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.Unauthorized("invalid credentials")
		}
		return nil, apperrors.Internal("failed to load user")
	}

	if err := s.passwords.Verify(input.Password, userRecord.PasswordHash); err != nil {
		return nil, apperrors.Unauthorized("invalid credentials")
	}

	response, appErr := s.issueTokens(ctx, userRecord)
	if appErr != nil {
		return nil, appErr
	}

	s.analytics.Track(ctx, "auth.logged_in", map[string]any{"user_id": userRecord.ID})

	return response, nil
}

func (s *Service) Refresh(ctx context.Context, input RefreshRequest) (*AuthResponse, *apperrors.AppError) {
	tokenHash := utils.HashToken(input.RefreshToken)
	session, err := s.repo.GetRefreshSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.Unauthorized("invalid refresh token")
		}
		return nil, apperrors.Internal("failed to load refresh session")
	}

	if session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
		return nil, apperrors.Unauthorized("refresh token is no longer valid")
	}

	if err := s.repo.RevokeRefreshSession(ctx, session.ID, time.Now().UTC()); err != nil {
		return nil, apperrors.Internal("failed to rotate refresh token")
	}

	userRecord, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, apperrors.Internal("failed to load user")
	}

	response, appErr := s.issueTokens(ctx, userRecord)
	if appErr != nil {
		return nil, appErr
	}

	return response, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*ProfileResponse, *apperrors.AppError) {
	userRecord, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.NotFound("user not found")
		}
		return nil, apperrors.Internal("failed to load user")
	}

	profile := profileFromUser(userRecord)
	return &profile, nil
}

func (s *Service) ForgotPassword(ctx context.Context, input ForgotPasswordRequest) *apperrors.AppError {
	userRecord, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return apperrors.Internal("failed to process password reset")
	}

	token, err := utils.RandomToken(32)
	if err != nil {
		return apperrors.Internal("failed to generate reset token")
	}

	tokenHash := utils.HashToken(token)
	expiresAt := time.Now().UTC().Add(s.resetTokenTTL)
	if err := s.repo.SetResetToken(ctx, userRecord.ID, tokenHash, expiresAt); err != nil {
		return apperrors.Internal("failed to persist reset token")
	}

	body := fmt.Sprintf("Use this token to reset your password: %s", token)
	if err := s.mailer.Send(ctx, userRecord.Email, "Reset your password", body); err != nil {
		return apperrors.Internal("failed to send reset email")
	}

	s.analytics.Track(ctx, "auth.password_reset_requested", map[string]any{"user_id": userRecord.ID})

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordRequest) *apperrors.AppError {
	tokenHash := utils.HashToken(input.Token)
	userRecord, err := s.repo.GetUserByResetTokenHash(ctx, tokenHash)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.Unauthorized("invalid reset token")
		}
		return apperrors.Internal("failed to load reset token")
	}

	if userRecord.ResetTokenExpires == nil || time.Now().After(*userRecord.ResetTokenExpires) {
		return apperrors.Unauthorized("reset token has expired")
	}

	passwordHash, err := s.passwords.Hash(input.NewPassword)
	if err != nil {
		return apperrors.Internal("failed to hash password")
	}

	if err := s.repo.ResetPassword(ctx, userRecord.ID, passwordHash); err != nil {
		return apperrors.Internal("failed to update password")
	}

	_ = s.repo.RevokeUserRefreshSessions(ctx, userRecord.ID, time.Now().UTC())
	s.analytics.Track(ctx, "auth.password_reset_completed", map[string]any{"user_id": userRecord.ID})

	return nil
}

func (s *Service) issueTokens(ctx context.Context, userRecord *UserRecord) (*AuthResponse, *apperrors.AppError) {
	accessToken, err := s.jwt.Generate(userRecord.ID, userRecord.Role, userRecord.Permissions)
	if err != nil {
		return nil, apperrors.Internal("failed to generate access token")
	}

	refreshToken, err := utils.RandomToken(32)
	if err != nil {
		return nil, apperrors.Internal("failed to generate refresh token")
	}

	_, err = s.repo.CreateRefreshSession(
		ctx,
		userRecord.ID,
		utils.HashToken(refreshToken),
		time.Now().UTC().Add(s.refreshTokenTTL),
	)
	if err != nil {
		return nil, apperrors.Internal("failed to persist refresh token")
	}

	profile := profileFromUser(userRecord)

	return &AuthResponse{
		User: profile,
		Tokens: TokenPairResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    s.tokenTTL,
		},
	}, nil
}

func profileFromUser(userRecord *UserRecord) ProfileResponse {
	return ProfileResponse{
		ID:          userRecord.ID,
		Email:       userRecord.Email,
		Role:        userRecord.Role,
		Permissions: userRecord.Permissions,
	}
}
