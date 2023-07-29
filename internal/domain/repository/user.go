package repository

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"

)

type UserRepository interface {
	Create(ctx context.Context, input model.User) error
	GetByСredentials(ctx context.Context, email, password string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserById(ctx context.Context, userID string) (model.User, error)
	GetUserByVerificationCode(ctx context.Context, hash string) (model.User, error)
	ChangeVerificationCode(ctx context.Context, email string, input model.UserVerificationPayload) error
	Verify(ctx context.Context, verificationCode string) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error)
	SetSession(ctx context.Context, userID string, session model.Session, lastVisitTime time.Time) error
	IsDuplicate(ctx context.Context, email string) (bool, error)
}
