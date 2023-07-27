package service

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/email"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"

)

type UserServiceI interface {
	SignUp(ctx context.Context, input types.UserSignUpDTO) (types.VerificationCodeDTO, error)
	SignIn(ctx context.Context, input types.UserSignInDTO) (types.Tokens, error)
	ResendVerificationCode(ctx context.Context, email string) (types.VerificationCodeDTO, error)
	RefreshTokens(ctx context.Context, refreshToken string) (types.Tokens, error)
	Verify(ctx context.Context, email string, verificationCode string) error
}

type Service struct {
	UserService UserServiceI
}

type Deps struct {
	Hasher              *hash.Hasher
	UserRepository      repository.UserRepository
	JWTManager          *manager.JWTManager
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	VerificationCodeTTL time.Duration
	Sender              email.Sender
	EmailConfig         config.EmailConfig
	CodeGenerator       *codegenerator.CodeGenerator
}

func NewService(deps *Deps) *Service {
	emailService := NewEmailService(deps.Sender, deps.EmailConfig)
	return &Service{
		UserService: NewUserService(deps.Hasher, deps.UserRepository, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService),
	}
}
