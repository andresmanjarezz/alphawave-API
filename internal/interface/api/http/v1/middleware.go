package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
	teamCtx             = "teamID"
	roleCtx             = "role"
)

func getUserId(c *gin.Context) (string, error) {
	return getIdByContext(c, userCtx)
}

func getTeamId(c *gin.Context) (string, error) {
	return getIdByContext(c, teamCtx)
}

func getRole(c *gin.Context) (string, error) {
	idFromCtx, ok := c.Get(roleCtx)
	if !ok {
		return "", errors.New("roleCtx not found")
	}

	id, ok := idFromCtx.(string)
	if !ok {
		return "", errors.New("roleCtx is of invalid type")
	}

	return id, nil
}

func getIdByContext(c *gin.Context, context string) (string, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return "", errors.New("userCtx not found")
	}

	id, ok := idFromCtx.(string)
	if !ok {
		return "", errors.New("userCtx is of invalid type")
	}

	return id, nil
}

func (h *HandlerV1) setTeamSessionFromCookie(c *gin.Context) {
	sessionDataCookie, err := c.Cookie("team_session")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	var sessionData teamSession

	if err := json.Unmarshal([]byte(sessionDataCookie), &sessionData); err != nil {
		newResponse(c, http.StatusInternalServerError, "error: error unmarshal data to json")
		return
	}

	c.Set(teamCtx, sessionData.TeamID)
	c.Set(roleCtx, sessionData.Roles)
}

func (h *HandlerV1) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

func (h *HandlerV1) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.JWTManager.ParseJWT(headerParts[1])
}

// func (h *HandlerV1) checkPermissions(permission string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		teamID, err := getTeamId(c)
// 		if err != nil {
// 			newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
// 			return
// 		}
// 		role, err := getRole(c)
// 		if err != nil {
// 			newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
// 			return
// 		}
// 		roles, err := h.service.RolesService.GetRolesByTeamId(c.Request.Context(), teamID)
// 		if err != nil {
// 			newResponse(c, http.StatusForbidden, errors.New("forbidden").Error())
// 			return
// 		}
// 		for _, item := range roles {
// 			if role == item.Role {
// 				switch
// 			}
// 		}
// 	}
// }
