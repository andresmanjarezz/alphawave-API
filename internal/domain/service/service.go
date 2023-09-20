package service

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
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
	UpdateUserInfo(ctx context.Context, userID string, input types.UpdateUserInfoDTO) error
	UpdateUserSettings(ctx context.Context, userID string, input types.UpdateUserSettingsDTO) error
	ResetPassword(ctx context.Context, email, token, tokenResult, password string) error
	VerifyForgotPasswordToken(ctx context.Context, email, token, tokenResult string) (types.ForgotPasswordPayloadDTO, error)
	ForgotPassword(ctx context.Context, email string) error
	ResendVerificationCode(ctx context.Context, email string) error
	RefreshTokens(ctx context.Context, refreshToken string) (types.Tokens, error)
	Verify(ctx context.Context, verificationCode string) (types.Tokens, error)
}

type MemberServiceI interface {
	MemberSignUp(ctx context.Context, token string, input types.MemberSignUpDTO) error
	GetMembersByQuery(ctx context.Context, teamID string, query types.GetMembersByQuery) ([]types.MemberDTO, error)
	GetMemberByTeamIdAndUserId(ctx context.Context, teamID string, userID string) (model.Member, error)
	UserInvite(ctx context.Context, teamID string, email string, role string) error
	AcceptInvite(ctx context.Context, token string) (string, error)
}

type TeamsServiceI interface {
	Create(ctx context.Context, userID string, input types.CreateTeamsDTO) error
	GetTeamByID(ctx context.Context, teamID string) (model.Team, error)
	GetTeamsByUser(ctx context.Context, userID string) ([]model.Team, error)
	UpdateTeamSettings(ctx context.Context, teamID string, input types.UpdateTeamSettingsDTO) error
}

type RolesServiceI interface {
	Create(ctx context.Context, teamID string) error
	GetRolesByTeamId(ctx context.Context, teamID string) ([]types.GetRoleDTO, error)
	UpdatePermissions(ctx context.Context, teamID string, input []types.UpdatePermissionsDTO) error
}

type TasksServiceI interface {
	Create(ctx context.Context, userID string, input types.TasksCreateDTO) error
	GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error)
	GetAll(ctx context.Context, userID string) ([]types.TaskDTO, error)
	UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error)
	ChangeStatus(ctx context.Context, userID, taskID, status string) error
	DeleteAll(ctx context.Context, userID string, status string) error
}

type AiChatServiceI interface {
	NewMessage(messages []types.Message) (types.Message, error)
}

type ProjectsServiceI interface {
}

type Service struct {
	UserService     UserServiceI
	MemberService   MemberServiceI
	TasksService    TasksServiceI
	RolesService    RolesServiceI
	ProjectsService ProjectsServiceI
	TeamsService    TeamsServiceI
	AiChatService   AiChatServiceI
}

type Deps struct {
	Hasher                 *hash.Hasher
	UserRepository         repository.UserRepository
	MemberRepository       repository.MemberRepository
	TasksRepository        repository.TasksRepository
	ProjectsRepository     repository.ProjectsRepository
	TeamsRepository        repository.TeamsRepository
	RolesRepository        repository.RolesRepository
	JWTManager             *manager.JWTManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeTTL    time.Duration
	Sender                 email.Sender
	EmailConfig            config.EmailConfig
	CodeGenerator          *codegenerator.CodeGenerator
	OpenAI                 openAI
	VerificationCodeLength int
	ApiUrl                 string
}

func NewService(deps *Deps) *Service {
	emailService := NewEmailService(deps.Sender, deps.EmailConfig)
	rolesService := NewRolesService(deps.RolesRepository)
	teamsService := NewTeamsService(deps.TeamsRepository, deps.UserRepository, deps.MemberRepository, *rolesService)
	userService := NewUserService(deps.Hasher, deps.UserRepository, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService, deps.VerificationCodeLength, deps.ApiUrl)
	return &Service{
		AiChatService:   NewAiChatService(deps.OpenAI),
		UserService:     userService,
		MemberService:   NewMemberService(deps.MemberRepository, deps.UserRepository, deps.CodeGenerator, teamsService, emailService, userService, deps.ApiUrl),
		TeamsService:    teamsService,
		RolesService:    rolesService,
		TasksService:    NewTasksService(deps.TasksRepository),
		ProjectsService: NewProjectsService(deps.ProjectsRepository),
	}
}
