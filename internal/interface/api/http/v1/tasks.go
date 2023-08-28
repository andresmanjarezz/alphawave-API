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
			authenticated.POST("/change-status", h.changeStatus)
			authenticated.POST("/update", h.updateByIdTask)
			authenticated.DELETE("/:status", h.deleteAll)
		}
	}
}

type CreateTaskInput struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

type UpdateTaskInput struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Order    int    `json:"order"`
}

type ChangeStatusInput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
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
		Title:    input.Title,
		Status:   input.Status,
		Priority: input.Priority,
		Order:    input.Order,
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

func (h *HandlerV1) updateByIdTask(c *gin.Context) {
	var input UpdateTaskInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	res, err := h.service.TasksService.UpdateById(c.Request.Context(), userID, types.UpdateTaskDTO{
		ID:       input.ID,
		Title:    input.Title,
		Status:   input.Status,
		Priority: input.Priority,
		Order:    input.Order,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.JSON(http.StatusFound, res)
}

func (h *HandlerV1) changeStatus(c *gin.Context) {
	var input ChangeStatusInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	if err := h.service.TasksService.ChangeStatus(c.Request.Context(), userID, input.ID, input.Status); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) deleteAll(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}
	status := c.Param("status")
	if status == "" {
		newResponse(c, http.StatusBadRequest, "status is empty")
		return
	}

	err = h.service.TasksService.DeleteAll(c.Request.Context(), userID, status)
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
