package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.GET("/verify/:code", h.userVerify)
		users.POST("/resend-verification", h.resendVerificationCode)
		users.GET("/auth/refresh", h.userRefresh)
		users.POST("/forgot-password", h.forgotPassword)
		users.GET("/forgot-password-verify", h.verifyForgotPasswordToken)
		users.POST("/reset-password", h.resetPassword)
		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.GET("/me", h.getUser)
			authenticated.POST("/change-password", h.changePassword)
			authenticated.PUT("/", h.updateUserInfo)
			authenticated.PUT("/settings", h.updateUserSettings)
			authenticated.POST("/logout", h.logOut)
		}
	}

}

type ChangePasswordInput struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type UserSignUpInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	JobTitle  string `json:"jobTitle"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserSignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordInput struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	TokenResult string `json:"tokenResult"`
	Password    string `json:"password"`
}

type UserVerifyInput struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verificationCode"`
}

type UpdateUserInfoInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	JobTitle  *string `json:"jobTitle"`
	Email     *string `json:"email"`
}

type UpdateUserSettingsInput struct {
	TimeZone   *string `json:"timeZone"`
	DateFormat *string `json:"dateFormat"`
	TimeFormat *string `json:"timeFormat"`
}

// type refreshTokenInput struct {
// 	Token string `json:"token" binding:"required"`
// }

type tokenResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	MattermostToken string `json:"mattermostToken"`
}

// type verifyResponse struct {
// 	Email                       string        `json:"email"`
// 	VerificationCodeExpiresTime time.Duration `json:"verificationCodeExpiresTime"`
// }

type EmailInput struct {
	Email string `json:"email"`
}

func (h *HandlerV1) signUp(c *gin.Context) {
	var input UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	err := h.service.UserService.SignUp(c.Request.Context(), types.UserSignUpDTO{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     input.Email,
		Password:  input.Password,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			newResponse(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectEmailFormat) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectPasswordFormat) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrIncorrectUserData) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *HandlerV1) getUser(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	user, err := h.service.UserService.GetUserById(c.Request.Context(), userID)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) resendVerificationCode(c *gin.Context) {
	var input EmailInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err := h.service.UserService.ResendVerificationCode(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) signIn(c *gin.Context) {
	var input UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	tokens, err := h.service.UserService.SignIn(c.Request.Context(), types.UserSignInDTO{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrUserNotVerifyed) {
			newResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, int(h.refreshTokenTTL.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:     tokens.AccessToken,
		RefreshToken:    tokens.RefreshToken,
		MattermostToken: tokens.MattermostToken,
	})
}

func (h *HandlerV1) userRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.service.UserService.RefreshTokens(c.Request.Context(), refreshToken)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	c.SetCookie("refresh_token", res.RefreshToken, int(h.refreshTokenTTL.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:     res.AccessToken,
		RefreshToken:    res.RefreshToken,
		MattermostToken: res.MattermostToken,
	})
}

func (h *HandlerV1) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return
	}
	tokens, err := h.service.UserService.Verify(c.Request.Context(), code)

	if err != nil {
		if errors.Is(err, apperrors.ErrIncorrectVerificationCode) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrUserAlreadyVerifyed) {
			newResponse(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, apperrors.ErrVerificationCodeExpired) {
			newResponse(c, http.StatusGone, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	c.SetCookie("refresh_token", tokens.RefreshToken, int(h.refreshTokenTTL.Seconds()), "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://%s/create-team?access_token=%s/mattermost_token=%s", h.frontEndUrl, tokens.AccessToken, tokens.MattermostToken))
	// c.String(http.StatusOK, "success")

}

func (h *HandlerV1) changePassword(c *gin.Context) {
	var input ChangePasswordInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	err = h.service.UserService.ChangePassword(c.Request.Context(), userID, input.NewPassword, input.OldPassword)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.String(http.StatusOK, "success")
}

func (h *HandlerV1) forgotPassword(c *gin.Context) {
	var input EmailInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err := h.service.UserService.ForgotPassword(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		if errors.Is(err, apperrors.ErrUserBlocked) {
			newResponse(c, http.StatusForbidden, apperrors.ErrUserBlocked.Error())
			return
		}
		logger.Error(err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) verifyForgotPasswordToken(c *gin.Context) {
	var (
		email       = c.Query("email")
		token       = c.Query("token")
		tokenResult = c.Query("result")
	)

	if email == "" || token == "" || tokenResult == "" {
		newResponse(c, http.StatusBadRequest, "url is incorrect")
		logger.Error("emial or token or result is empty")
		return
	}

	tokenPayload, err := h.service.UserService.VerifyForgotPasswordToken(c.Request.Context(), email, token, tokenResult)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, "url is incorrect")
			logger.Error(err)
			return
		}
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, apperrors.ErrInternalServerError.Error())
		return
	}

	url := fmt.Sprintf("http://%s/reset-password?email=%s&token=%s&result=%s", h.frontEndUrl, tokenPayload.Email, tokenPayload.Token, tokenPayload.ResultToken)
	c.Redirect(http.StatusFound, url)
}

func (h *HandlerV1) resetPassword(c *gin.Context) {
	var input ResetPasswordInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err := h.service.UserService.ResetPassword(c.Request.Context(), input.Email, input.Token, input.TokenResult, input.Password)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) updateUserInfo(c *gin.Context) {

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	var input UpdateUserInfoInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	if err := h.service.UserService.UpdateUserInfo(c.Request.Context(), userID, types.UpdateUserInfoDTO{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     input.Email,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) updateUserSettings(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrDocumentNotFound.Error())
		return
	}
	var input UpdateUserSettingsInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}
	if err := h.service.UserService.UpdateUserSettings(c.Request.Context(), userID, types.UpdateUserSettingsDTO{
		TimeZone:   input.TimeZone,
		DateFormat: input.DateFormat,
		TimeFormat: input.TimeFormat,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) logOut(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrDocumentNotFound.Error())
		return
	}

	err = h.service.UserService.LogOut(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrUserNotFound.Error())
		return
	}
	c.Status(http.StatusOK)
}
