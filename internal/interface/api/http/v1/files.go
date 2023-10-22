package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 2 << 30
)

// type contentRange struct {
// 	rangeStart int64
// 	rangeEnd   int64
// 	fileSize   int64
// }

type CreateFolderInput struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type FileMetadataResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	Extension string    `json:"extension"`
	File      http.File `json:"file"`
}

func (h *HandlerV1) initFilesRoutes(api *gin.RouterGroup) {
	files := api.Group("/files")
	{
		authenticated := files.Group("/", h.userIdentity)
		{
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.GET("/url/:id", h.getFilePresignedURL)
				teamSession.GET("/:id", h.getFile)
				teamSession.POST("/", h.createFile)
				teamSession.POST("/folder", h.createFolder)
				teamSession.DELETE("/:id", h.deleteFile)
			}
		}

	}
}

func (h *HandlerV1) createFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	location := c.PostForm("location")

	fileReader, err := file.Open()

	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, fmt.Errorf("failed to open file").Error())
		return
	}

	defer fileReader.Close()

	fileData, err := ioutil.ReadAll(fileReader)

	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	reader := bytes.NewReader(fileData)

	fileName := filepath.Ext(file.Filename)

	if err := h.service.FilesService.Create(c.Request.Context(), reader, types.CreateFileDTO{
		TeamID:    teamID,
		FileName:  fileName,
		Location:  location,
		Extension: filepath.Ext(file.Filename),
		Size:      len(fileData),
	}); err != nil {

		logger.Errorf("failed to upload file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to upload file").Error())
		return
	}

	c.Status(http.StatusCreated)

}

func (h *HandlerV1) createFolder(c *gin.Context) {

	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	var input CreateFolderInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))
		return
	}

	if err := h.service.FilesService.CreateFolder(c.Request.Context(), types.CreateFolderDTO{
		TeamID:     teamID,
		FolderName: input.Name,
		Location:   input.Location,
	}); err != nil {
		logger.Errorf("failed to create folder. err: %v", err)
		newResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create folder. folder name: %s", input.Name).Error())
		return
	}
	c.JSON(http.StatusCreated, fmt.Sprintf("%s/%s", input.Location, input.Name))
}

func (h *HandlerV1) getFilePresignedURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Warnf("failed to get file id. file id is empty")
		newResponse(c, http.StatusBadRequest, "id is empty")
		return
	}
	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	url, err := h.service.FilesService.GetFilePresignedURL(c.Request.Context(), teamID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to get file presigned url. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		logger.Errorf("failed to get file presigned url. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("https://%s", url))
}

func (h *HandlerV1) getFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Warnf("failed to get file id. file id is empty")
		newResponse(c, http.StatusBadRequest, "id is empty")
		return
	}
	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id for get file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	file, err := h.service.FilesService.GetFile(c.Request.Context(), teamID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to get file. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		logger.Errorf("failed to get file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s; fileID=%s", file.Name, file.ID))
	c.Header("Content-Length", strconv.FormatInt(int64(file.Size), 10))

	c.Data(http.StatusOK, "application/octet-stream", *file.File)
}

// func (h *HandlerV1) getFolderContent(c *gin.Context) {

// 	objectName := c.Param("object")

// 	if objectName == "" {
// 		newResponse(c, http.StatusBadRequest, "object name is empty")
// 		return
// 	}

// 	var prefix string

// 	if objectName == "" {
// 		prefix = ""
// 	} else {
// 		prefix = fmt.Sprintf("%s/", objectName)
// 	}

// }

func (h *HandlerV1) deleteFile(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")
		return
	}

	teamID, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	if err := h.service.FilesService.Delete(c.Request.Context(), teamID, id); err != nil {

		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to delete file. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())
			return
		}
		logger.Errorf("failed to delete file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		return
	}

	c.Status(http.StatusNoContent)

}
