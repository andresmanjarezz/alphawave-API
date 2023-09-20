package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/service"
	openai "github.com/Coke15/AlphaWave-BackEnd/internal/infrastructure/ai/openAI"
	httpRoutes "github.com/Coke15/AlphaWave-BackEnd/internal/interface/api/http"
	"github.com/Coke15/AlphaWave-BackEnd/internal/repository"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/db/mongodb"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/email/smtp"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
)

const configDir = "configs"

func Run() {

	cfg, err := config.Init(configDir)

	if err != nil {
		logger.Errorf("error parse config. err: %v", err)
	}

	// -----
	hasher := hash.NewHasher(cfg.Auth.PasswordSalt)

	mongoClient, err := mongodb.NewConnection(cfg.MongoDB.Url, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		logger.Errorf("failed to create new mongo client. err: %v", err)
	}

	JWTManager, err := manager.NewJWTManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)
		return
	}
	codeGenerator := codegenerator.NewCodeGenerator()

	openAI := openai.NewOpenAiAPI(cfg.OpenAI.Token, cfg.OpenAI.Url)
	// -----

	mongodb := mongoClient.Database(cfg.MongoDB.DBName)
	repository := repository.NewRepository(mongodb)
	service := service.NewService(&service.Deps{
		UserRepository:         repository.User,
		TasksRepository:        repository.Tasks,
		ProjectsRepository:     repository.Projects,
		TeamsRepository:        repository.Teams,
		RolesRepository:        repository.Roles,
		MemberRepository:       repository.Members,
		Hasher:                 hasher,
		JWTManager:             JWTManager,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		VerificationCodeTTL:    cfg.Auth.VerificationCodeTTL,
		Sender:                 emailSender,
		EmailConfig:            cfg.Email,
		CodeGenerator:          codeGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		ApiUrl:                 cfg.HTTP.Host,
		OpenAI:                 openAI,
	})
	handler := httpRoutes.NewHandler(service, JWTManager, cfg.Auth.JWT.RefreshTokenTTL, cfg.FrontEndUrl)

	srv := NewServer(cfg, handler.InitRoutes(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	if err := srv.Shotdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %x", err)
	}
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Errorf("error disconnect to mongoClient. err: %v", err)
	}
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes << 20,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	port := strings.Replace(s.httpServer.Addr, ":", "", 1)

	logger.Infof("Server has ben started on port: %s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shotdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
