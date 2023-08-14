package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initTasksRoutes(api *gin.RouterGroup) {
	tasks := api.Group("/tasks")
	{
		authenticated := tasks.Group("/", h.userIdentity)
		{
			authenticated.POST("/create", h.createTask)
			authenticated.GET("/:id", h.getByIdTask)
			authenticated.GET("/", h.getAllTasks)
		}
	}
}

type CreateTaskInput struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

func (h *HandlerV1) createTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	var input CreateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	err = h.service.TasksService.Create(c.Request.Context(), userID, types.TasksCreateDTO{
		Title: input.Title,
		Order: input.Order,
	})

	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.Status(http.StatusCreated)

}

func (h *HandlerV1) getByIdTask(c *gin.Context) {
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

	task, err := h.service.TasksService.GetById(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *HandlerV1) getAllTasks(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	tasks, err := h.service.TasksService.GetAll(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
