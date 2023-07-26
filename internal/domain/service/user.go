package service

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
)

type UserService struct {
	hasher              *hash.Hasher
	repository          repository.UserRepository
	JWTManager          *manager.JWTManager
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	VerificationCodeTTL time.Duration
	codeGenerator       *codegenerator.CodeGenerator
	emailService        *EmailService
}

func NewUserService(hasher *hash.Hasher, repository repository.UserRepository, JWTManager *manager.JWTManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, verificationCodeTTL time.Duration, codeGenerator *codegenerator.CodeGenerator, emailService *EmailService) *UserService {
	return &UserService{
		hasher:              hasher,
		repository:          repository,
		JWTManager:          JWTManager,
		AccessTokenTTL:      accessTokenTTL,
		RefreshTokenTTL:     refreshTokenTTL,
		emailService:        emailService,
		VerificationCodeTTL: verificationCodeTTL,
		codeGenerator:       codeGenerator,
	}
}

func (s *UserService) SignUp(ctx context.Context, input types.UserSignUpDTO) (types.VerificationCodeDTO, error) {
	if err := validateCredentials(input.Email, input.Password); err != nil {
		return types.VerificationCodeDTO{}, err
	}
	if err := validateUserData(input.FirstName, input.LastName, input.JobTitle); err != nil {
		return types.VerificationCodeDTO{}, err
	}
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return types.VerificationCodeDTO{}, err
	}
	verificationCode := s.codeGenerator.GenerateUniqueCode()

	user := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     input.Email,
		Password:  passwordHash,
		Verification: model.UserVerificationPayload{
			VerificationCode:            verificationCode,
			VerificationCodeExpiresTime: time.Now().Add(s.VerificationCodeTTL),
		},
		RegisteredTime: time.Now(),
		LastVisitTime:  time.Now(),
	}

	isDuplicate, err := s.repository.IsDuplicate(ctx, input.Email)
	if err != nil {
		return types.VerificationCodeDTO{}, err

	}
	if isDuplicate {
		return types.VerificationCodeDTO{}, apperrors.ErrUserAlreadyExists
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return types.VerificationCodeDTO{}, err
	}
	go func() {
		err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
			Name:             input.FirstName,
			Email:            input.Email,
			VerificationCode: verificationCode,
		})
		logger.Error(err)
	}()

	return types.VerificationCodeDTO{
		Email:                       input.Email,
		VerificationCodeExpiresTime: s.VerificationCodeTTL / time.Second,
	}, nil
}

func (s *UserService) SignIn(ctx context.Context, input types.UserSignInDTO) (types.Tokens, error) {
	err := validateCredentials(input.Email, input.Password)
	if err != nil {
		return types.Tokens{}, err
	}
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return types.Tokens{}, err
	}

	user, err := s.repository.GetByСredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.Tokens{}, err
		}
		return types.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) LogOut(ctx context.Context, userID string) {

}

func (s *UserService) Verify(ctx context.Context, email string, verificationCode string) error {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}
	if user.Verification.Verified == true {
		return apperrors.ErrUserAlreadyVerifyed
	}
	if user.Verification.VerificationCode != verificationCode {
		return apperrors.ErrIncorrectVerificationCode
	}
	if user.Verification.VerificationCodeExpiresTime.UTC().Unix() < time.Now().UTC().Unix() {
		return apperrors.ErrVerificationCodeExpired
	}

	return s.repository.Verify(ctx, user.ID, verificationCode)
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (types.Tokens, error) {
	user, err := s.repository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.Tokens{}, err
		}
		return types.Tokens{}, err
	}
	if user.Blocked {
		return types.Tokens{}, apperrors.ErrUserBlocked
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) createSession(ctx context.Context, userID string) (types.Tokens, error) {

	accessToken, err := s.JWTManager.NewJWT(userID, s.AccessTokenTTL)
	if err != nil {
		return types.Tokens{}, err
	}
	refreshToken, err := s.JWTManager.NewRefreshToken()
	if err != nil {
		return types.Tokens{}, err
	}
	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	session := model.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresTime:  time.Now().Add(s.RefreshTokenTTL),
	}

	err = s.repository.SetSession(ctx, userID, session, time.Now())
	return types.Tokens(tokens), err
}
