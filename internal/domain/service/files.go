package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
)

type storageProvider interface {
	UploadFile(ctx context.Context, bucketName, objectName, fileName string, fileSize int64, reader io.Reader) error
	CreateFolder(ctx context.Context, bucketName, objectName string) error
	GetFilePresignedURL(ctx context.Context, bucketName, fileName string, expiresTime time.Duration) (string, error)
	GetFile(ctx context.Context, bucketName, fileName string) (*[]byte, error)
	DeleteFile(ctx context.Context, bucketName, fileName string) error
}

type FilesService struct {
	storageProvider storageProvider
	repository      repository.FilesRepository
	codeGenerator   *codegenerator.CodeGenerator
}

const defaultFileExpiresTime = time.Second * 60 * 60 * 24

func NewFilesService(storageProvider storageProvider, repository repository.FilesRepository, codeGenerator *codegenerator.CodeGenerator) *FilesService {
	return &FilesService{
		storageProvider: storageProvider,
		repository:      repository,
		codeGenerator:   codeGenerator,
	}
}

func (s *FilesService) Create(ctx context.Context, reader io.Reader, input types.CreateFileDTO) error {
	uuid := s.codeGenerator.GenerateUUID()

	fileName := fmt.Sprintf("%s/%s%s", input.Location, uuid, input.Extension)

	fileID, err := s.repository.Create(ctx, model.File{
		TeamId:    input.TeamID,
		Name:      fmt.Sprintf("%s/%s", input.Location, input.FileName),
		FilePath:  fileName,
		Key:       uuid,
		Type:      "file",
		Size:      input.Size,
		Extension: input.Extension[1:],
	})

	if err != nil {
		return err
	}

	if err := s.storageProvider.UploadFile(ctx, model.BUCKET_STORAGE, fileName, input.FileName, int64(input.Size), reader); err != nil {

		if err := s.repository.Delete(ctx, input.TeamID, fileID); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (s *FilesService) CreateFolder(ctx context.Context, input types.CreateFolderDTO) error {
	var folderName string

	if input.Location == "" {
		folderName = input.FolderName + "/"
	} else {
		folderName = input.Location + "/" + input.FolderName + "/"
	}

	if err := s.storageProvider.CreateFolder(ctx, model.BUCKET_STORAGE, folderName); err != nil {
		return err
	}

	return nil
}

// func (s *FilesService) GetFolderContent()

func (s *FilesService) GetFilePresignedURL(ctx context.Context, teamID, fileID string) (string, error) {
	file, err := s.repository.GetFileById(ctx, teamID, fileID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return "", apperrors.ErrDocumentNotFound
		}
		return "", err
	}

	url, err := s.storageProvider.GetFilePresignedURL(ctx, model.BUCKET_STORAGE, file.FilePath, defaultFileExpiresTime)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *FilesService) GetFile(ctx context.Context, teamID, fileID string) (*types.GetFileDTO, error) {
	res, err := s.repository.GetFileById(ctx, teamID, fileID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return nil, apperrors.ErrDocumentNotFound
		}
		return nil, err
	}

	file, err := s.storageProvider.GetFile(ctx, model.BUCKET_STORAGE, res.FilePath)

	if err != nil {

		return nil, err
	}

	return &types.GetFileDTO{
		ID:        res.ID,
		Name:      res.Name,
		Size:      res.Size,
		Extension: res.Extension,
		File:      file,
	}, nil

}

func (s *FilesService) Delete(ctx context.Context, teamId, fileId string) error {
	res, err := s.repository.GetFileById(ctx, teamId, fileId)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}

		return err
	}

	if err := s.storageProvider.DeleteFile(ctx, model.BUCKET_STORAGE, res.FilePath); err != nil {
		return err
	}

	if err := s.repository.Delete(ctx, teamId, res.ID); err != nil {
		return err
	}
	return nil
}
