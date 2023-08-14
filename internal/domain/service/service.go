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
	SignUp(ctx context.Context, input types.UserSignUpDTO) error
	SignIn(ctx context.Context, input types.UserSignInDTO) (types.Tokens, error)
	GetUserById(ctx context.Context, userID string) (types.UserDTO, error)
	ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error
	ResetPassword(ctx context.Context, email, token, tokenResult, password string) error
	VerifyForgotPasswordToken(ctx context.Context, email, token, tokenResult string) (types.ForgotPasswordPayloadDTO, error)
	ForgotPassword(ctx context.Context, email string) error
	ResendVerificationCode(ctx context.Context, email string) error
	RefreshTokens(ctx context.Context, refreshToken string) (types.Tokens, error)
	Verify(ctx context.Context, verificationCode string) error
}

type TasksServiceI interface {
	Create(ctx context.Context, userID string, input types.TasksCreateDTO) error
	GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error)
	GetAll(ctx context.Context, userID string) ([]types.TaskDTO, error)
	UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error)
}

type Service struct {
	UserService  UserServiceI
	TasksService TasksServiceI
}

type Deps struct {
	Hasher                 *hash.Hasher
	UserRepository         repository.UserRepository
	TasksRepository        repository.TasksRepository
	JWTManager             *manager.JWTManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeTTL    time.Duration
	Sender                 email.Sender
	EmailConfig            config.EmailConfig
	CodeGenerator          *codegenerator.CodeGenerator
	VerificationCodeLength int
	ApiUrl                 string
}

func NewService(deps *Deps) *Service {
	emailService := NewEmailService(deps.Sender, deps.EmailConfig)
	return &Service{
		UserService:  NewUserService(deps.Hasher, deps.UserRepository, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService, deps.VerificationCodeLength, deps.ApiUrl),
		TasksService: NewTasksService(deps.TasksRepository),
	}
}
