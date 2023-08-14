package service

import (
	"context"
	"errors"
	"fmt"
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
	hasher                 *hash.Hasher
	repository             repository.UserRepository
	JWTManager             *manager.JWTManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeTTL    time.Duration
	VerificationCodeLength int
	ApiUrl                 string
	codeGenerator          *codegenerator.CodeGenerator
	emailService           *EmailService
}

func NewUserService(hasher *hash.Hasher, repository repository.UserRepository, JWTManager *manager.JWTManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, verificationCodeTTL time.Duration, codeGenerator *codegenerator.CodeGenerator, emailService *EmailService, verificationCodeLength int, apiUrl string) *UserService {
	return &UserService{
		hasher:                 hasher,
		repository:             repository,
		JWTManager:             JWTManager,
		AccessTokenTTL:         accessTokenTTL,
		RefreshTokenTTL:        refreshTokenTTL,
		emailService:           emailService,
		VerificationCodeTTL:    verificationCodeTTL,
		codeGenerator:          codeGenerator,
		VerificationCodeLength: verificationCodeLength,
		ApiUrl:                 apiUrl,
	}
}

func (s *UserService) SignUp(ctx context.Context, input types.UserSignUpDTO) error {
	if err := validateCredentials(input.Email, input.Password); err != nil {
		return err
	}
	if err := validateUserData(input.FirstName, input.LastName, input.JobTitle); err != nil {
		return err
	}
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	verificationCodeHash, err := s.hasher.Hash(input.Email)
	if err != nil {
		return err
	}
	verificationCode := fmt.Sprintf("%s%s", s.codeGenerator.RandomSecret(s.VerificationCodeLength), verificationCodeHash)

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
		return err

	}
	if isDuplicate {
		return apperrors.ErrUserAlreadyExists
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return err
	}
	go func() {
		err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
			Name:  input.FirstName,
			Email: input.Email,
			URL:   s.ApiUrl + "/api/v1/users/verify/" + verificationCode,
		})
		logger.Error(err)
	}()

	return nil
}

func (s *UserService) SignIn(ctx context.Context, input types.UserSignInDTO) (types.Tokens, error) {

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

	if !user.Verification.Verified {
		return types.Tokens{}, apperrors.ErrUserNotVerifyed
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) LogOut(ctx context.Context, userID string) {

}

func (s *UserService) EnableTwoFactorAuth(ctx context.Context, userID string) error {

	return nil
}

func (s *UserService) GetUserById(ctx context.Context, userID string) (types.UserDTO, error) {
	res, err := s.repository.GetUserById(ctx, userID)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.UserDTO{}, err
		}
		return types.UserDTO{}, err
	}

	user := types.UserDTO{
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		JobTitle:       res.JobTitle,
		Email:          res.Email,
		LastVisitTime:  res.LastVisitTime,
		RegisteredTime: res.RegisteredTime,
		Verification:   res.Verification.Verified,
		Blocked:        res.Blocked,
	}

	return user, nil
}

func (s *UserService) Verify(ctx context.Context, verificationCode string) error {
	user, err := s.repository.GetUserByVerificationCode(ctx, verificationCode)
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

	return s.repository.Verify(ctx, verificationCode)
}

func (s *UserService) ResendVerificationCode(ctx context.Context, email string) error {

	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	verificationCodeHash, err := s.hasher.Hash(email)
	if err != nil {
		return err
	}
	verificationCode := fmt.Sprintf("%s%s", s.codeGenerator.RandomSecret(s.VerificationCodeLength), verificationCodeHash)

	verificationPayload := model.UserVerificationPayload{
		VerificationCode:            verificationCode,
		VerificationCodeExpiresTime: time.Now().Add(s.VerificationCodeTTL),
	}
	err = s.repository.ChangeVerificationCode(ctx, email, verificationPayload)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	go func() {
		err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
			Name:  user.FirstName,
			Email: user.Email,
			URL:   s.ApiUrl + "/api/v1/users/verify/" + verificationCode,
		})
		logger.Error(err)
	}()
	return nil
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

func (s *UserService) ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error {

	passwordHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	oldPasswordHash, err := s.hasher.Hash(oldPassword)
	if err != nil {
		return err
	}

	err = s.repository.ChangePassword(ctx, userID, passwordHash, oldPasswordHash)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	if user.Blocked {
		return apperrors.ErrUserBlocked
	}
	tokenHash, err := s.hasher.Hash(user.Email)
	if err != nil {
		return err
	}

	result := fmt.Sprintf("%s.%s", s.codeGenerator.RandomSecret(30), tokenHash)

	tokenExpiresTime := time.Now().Add(time.Hour * 1)

	err = s.repository.SetForgotPassword(ctx, user.Email, model.ForgotPasswordPayload{
		Token:            tokenHash,
		ResultToken:      result,
		TokenExpiresTime: tokenExpiresTime,
	})

	if err != nil {
		return err
	}
	if err := s.emailService.SendUserForgotPassword(ForgotPasswordInput{
		Email:            user.Email,
		TokenExpiresTime: 1,
		URL:              s.ApiUrl + fmt.Sprintf("/api/v1/users/forgot-password-verify?email=%s&token=%s&result=%s", user.Email, tokenHash, result),
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyForgotPasswordToken(ctx context.Context, email, token, tokenResult string) (types.ForgotPasswordPayloadDTO, error) {

	user, err := s.repository.GetByForgotPasswordToken(ctx, token, tokenResult)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {

			return types.ForgotPasswordPayloadDTO{}, apperrors.ErrUserNotFound
		}
		return types.ForgotPasswordPayloadDTO{}, err
	}

	if user.Email != email {

		return types.ForgotPasswordPayloadDTO{}, apperrors.ErrUserNotFound
	}

	return types.ForgotPasswordPayloadDTO{
		Email:       user.Email,
		Token:       token,
		ResultToken: tokenResult,
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, email, token, tokenResult, password string) error {
	user, err := s.repository.GetByForgotPasswordToken(ctx, token, tokenResult)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	if user.Email != email {
		return apperrors.ErrUserNotFound
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	err = s.repository.ResetPassword(ctx, token, user.Email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
