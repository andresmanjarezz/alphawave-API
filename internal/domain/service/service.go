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

type MemberServiceI interface {
	GetMembers(ctx context.Context, teamID string, query types.GetUsersByQuery) ([]types.MemberDTO, error)
}

type TeamsServiceI interface {
	Create(ctx context.Context, userEmail string, input types.CreateTeamsDTO) error
}

type TasksServiceI interface {
	Create(ctx context.Context, userID string, input types.TasksCreateDTO) error
	GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error)
	GetAll(ctx context.Context, userID string) ([]types.TaskDTO, error)
	UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error)
	ChangeStatus(ctx context.Context, userID, taskID, status string) error
	DeleteAll(ctx context.Context, userID string, status string) error
}

type ProjectsServiceI interface {
}

type Service struct {
	UserService     UserServiceI
	MemberService   MemberServiceI
	TasksService    TasksServiceI
	ProjectsService ProjectsServiceI
	TeamsService    TeamsServiceI
}

type Deps struct {
	Hasher                 *hash.Hasher
	UserRepository         repository.UserRepository
	MemberRepository       repository.MemberRepository
	TasksRepository        repository.TasksRepository
	ProjectsRepository     repository.ProjectsRepository
	TeamsRepository        repository.TeamsRepository
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
		UserService:     NewUserService(deps.Hasher, deps.UserRepository, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService, deps.VerificationCodeLength, deps.ApiUrl),
		MemberService:   NewMemberService(deps.MemberRepository, deps.UserRepository),
		TeamsService:    NewTeamsService(deps.TeamsRepository, deps.UserRepository),
		TasksService:    NewTasksService(deps.TasksRepository),
		ProjectsService: NewProjectsService(deps.ProjectsRepository),
	}
}
