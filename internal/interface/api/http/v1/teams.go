package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initTeamsRoutes(api *gin.RouterGroup) {
	teams := api.Group("/teams")
	{
		authenticated := teams.Group("/", h.userIdentity)
		{
			authenticated.POST("/create", h.createTeam)
			authenticated.GET("/set-session/:id", h.setSession)
			authenticated.GET("/", h.getTeams)
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.PUT("/settings", h.updateSettings)
				roles := teamSession.Group("/roles")
				{
					roles.GET("/", h.getRoles)
					roles.PUT("/", h.UpdatePermissions)
				}
			}
		}

	}

}

type createTeamInput struct {
	TeamName string `json:"teamName"`
	JobTitle string `json:"jobTitle"`
}

type updateTeamSettingsInput struct {
	UserActivityIndicator *bool   `json:"userActivityIndicator"`
	DisplayLinkPreview    *bool   `json:"displayLinkPreview"`
	DisplayFilePreview    *bool   `json:"displayFilePreview"`
	EnableGifs            *bool   `json:"enableGifs"`
	ShowWeekends          *bool   `json:"showWeekends"`
	FirstDayOfWeek        *string `json:"firstDayOfWeek"`
}

type teamSession struct {
	TeamID string
	Roles  []string
}

type updatePermissionsInput struct {
	Role        string            `json:"role"`
	Permissions model.Permissions `jons:"permissions"`
}

func (h *HandlerV1) createTeam(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	var input createTeamInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err = h.service.TeamsService.Create(c.Request.Context(), userID, types.CreateTeamsDTO{
		TeamName: input.TeamName,
		JobTitle: input.JobTitle,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *HandlerV1) setSession(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")
		return
	}

	team, err := h.service.TeamsService.GetTeamByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, "team not foud")
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	sessionData := teamSession{
		TeamID: team.ID,
		Roles:  []string{},
	}

	fmt.Printf("userID: %s, ownerID: %s", userID, team.OwnerID)
	if team.OwnerID != userID {
		member, err := h.service.MemberService.GetMemberByTeamIdAndUserId(c.Request.Context(), team.ID, userID)
		if err != nil {
			if errors.Is(err, apperrors.ErrMemberNotFound) {
				newResponse(c, http.StatusNotFound, apperrors.ErrMemberNotFound.Error())
				return
			}
			newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
			return
		}
		sessionData.Roles = member.Roles
	}

	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "error: error marshal data to json")
		return
	}
	c.SetCookie("team_session", string(sessionDataJson), 0, "/", "", false, true)
}

func (h *HandlerV1) getTeams(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	teams, err := h.service.TeamsService.GetTeamsByUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *HandlerV1) updateSettings(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	var input updateTeamSettingsInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err = h.service.TeamsService.UpdateTeamSettings(c.Request.Context(), id, types.UpdateTeamSettingsDTO{
		UserActivityIndicator: input.UserActivityIndicator,
		DisplayLinkPreview:    input.DisplayLinkPreview,
		DisplayFilePreview:    input.DisplayFilePreview,
		EnableGifs:            input.EnableGifs,
		ShowWeekends:          input.ShowWeekends,
		FirstDayOfWeek:        input.FirstDayOfWeek,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) getRoles(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	roles, err := h.service.RolesService.GetRolesByTeamId(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *HandlerV1) UpdatePermissions(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	var input []updatePermissionsInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	roles := make([]types.UpdatePermissionsDTO, len(input))

	for i := range roles {
		roles[i] = types.UpdatePermissionsDTO{
			Role:        input[i].Role,
			Permissions: input[i].Permissions,
		}
	}
	err = h.service.RolesService.UpdatePermissions(c.Request.Context(), id, roles)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
