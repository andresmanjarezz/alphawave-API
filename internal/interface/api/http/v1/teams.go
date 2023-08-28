package v1

import (
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initTeamsRoutes(api *gin.RouterGroup) {
	teams := api.Group("/teams")
	{
		teams.POST("/create", h.createTeam)

	}

}

type createTeamInput struct {
	UserEmail string `json:"userEmail"`
	TeamName  string `json:"teamName"`
	JobTitle  string `json:"jobTitle"`
}

func (h *HandlerV1) createTeam(c *gin.Context) {
	var input createTeamInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err := h.service.TeamsService.Create(c.Request.Context(), input.UserEmail, types.CreateTeamsDTO{
		TeamName: input.TeamName,
		JobTitle: input.JobTitle,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusCreated)
}
