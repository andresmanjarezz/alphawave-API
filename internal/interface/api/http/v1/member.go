package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initMembersRoutes(api *gin.RouterGroup) {
	teams := api.Group("/members")
	{
		teams.GET("/", h.getMembers)
	}
}

func (h *HandlerV1) getMembers(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

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

	members, err := h.service.MemberService.GetMembers(c.Request.Context(), id, types.GetUsersByQuery{
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
