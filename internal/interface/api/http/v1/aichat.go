package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initAiChatRoutes(api *gin.RouterGroup) {
	ai := api.Group("/ai")
	{
		ai.POST("/new-message", h.newMessage)
	}
}

func (h *HandlerV1) newMessage(c *gin.Context) {
	var input []types.Message

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	message, err := h.service.AiChatService.NewMessage(input)

	if err != nil {
		newResponse(c, http.StatusBadGateway, errors.New("error gateway").Error())
		return
	}

	c.JSON(http.StatusOK, message)
}
