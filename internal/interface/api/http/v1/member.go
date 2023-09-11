package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initMembersRoutes(api *gin.RouterGroup) {
	members := api.Group("/members")
	{
		members.GET("/accept_invite/:token", h.acceptInvite)
		members.POST("/sign-up", h.memberSignUp)
		authenticated := members.Group("/", h.userIdentity, h.setTeamSessionFromCookie)
		{
			authenticated.GET("/", h.getMembers)
			authenticated.POST("/invite", h.userInvite)
		}

	}
}

type UserInviteInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type MemberSignUpInput struct {
	Token     string `json:"token"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	JobTitle  string `json:"jobTitle"`
	Password  string `json:"password"`
}

func (h *HandlerV1) memberSignUp(c *gin.Context) {
	var input MemberSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err := h.service.MemberService.MemberSignUp(c.Request.Context(), input.Token, types.MemberSignUpDTO{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Password:  input.Password,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			newResponse(c, http.StatusConflict, apperrors.ErrUserAlreadyExists.Error())
			return
		}
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrMemberNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *HandlerV1) getMembers(c *gin.Context) {

	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	skip, err := strconv.Atoi(c.Query("skip"))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	members, err := h.service.MemberService.GetMembersByQuery(c.Request.Context(), id, types.GetMembersByQuery{
		PaginationQuery: types.PaginationQuery{
			Skip:  int64(skip),
			Limit: int64(limit),
		},
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *HandlerV1) userInvite(c *gin.Context) {
	var input UserInviteInput

	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err = h.service.MemberService.UserInvite(c.Request.Context(), id, input.Email, input.Role)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) acceptInvite(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		newResponse(c, http.StatusBadRequest, "token is empty")
		return
	}

	token, err := h.service.MemberService.AcceptInvite(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://%s/members/sign-up/%s", h.frontEndUrl, token))
			return
		}
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrMemberNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusOK)
}
