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
		UserID:   userID,
		Title:    input.Title,
		Status:   input.Status,
		Priority: input.Priority,
		Order:    input.Order,
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
		ID:       task.ID,
		Title:    task.Title,
		Status:   task.Status,
		Priority: task.Priority,
		Order:    task.Order,
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
			ID:       tasksIn[i].ID,
			Title:    tasksIn[i].Title,
			Status:   tasksIn[i].Status,
			Priority: tasksIn[i].Priority,
			Order:    tasksIn[i].Order,
		}
	}

	return tasks, nil
}

func (s *TasksService) UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error) {
	task, err := s.repository.UpdateById(ctx, userID, model.Task{
		ID:       input.ID,
		Title:    input.Title,
		Priority: input.Priority,
		Status:   input.Status,
		Order:    input.Order,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.TaskDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.TaskDTO{}, err
	}
	return types.TaskDTO{
		ID:       task.ID,
		Title:    task.Title,
		Priority: task.Priority,
		Status:   task.Status,
		Order:    task.Order,
	}, err
}
