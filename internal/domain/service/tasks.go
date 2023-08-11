package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type TasksService struct {
	repository repository.TasksRepository
}

func NewTasksService(repository repository.TasksRepository) *TasksService {
	return &TasksService{
		repository: repository,
	}
}

func (s *TasksService) Create(ctx context.Context, userID string, input types.TasksCreateDTO) error {

	task := model.Task{
		UserID: userID,
		Title:  input.Title,
		Order:  input.Order,
	}
	err := s.repository.Create(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TasksService) GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error) {
	task, err := s.repository.GetById(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.TaskDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.TaskDTO{}, err
	}

	return types.TaskDTO{
		ID:    task.ID,
		Title: task.Title,
		Order: task.Order,
	}, nil
}

func (s *TasksService) GetAll(ctx context.Context, userID string) ([]types.TaskDTO, error) {
	tasksIn, err := s.repository.GetAll(ctx, userID)

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.TaskDTO{}, err
		}
		return []types.TaskDTO{}, err
	}

	tasks := make([]types.TaskDTO, len(tasksIn))

	for i := range tasksIn {
		tasks[i] = types.TaskDTO{
			ID:    tasksIn[i].ID,
			Title: tasksIn[i].Title,
			Order: tasksIn[i].Order,
		}
	}

	return tasks, nil
}
