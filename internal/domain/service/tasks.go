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

const (
	STATUS_DELITE string = "del"
	STATUS_DONE   string = "done"
	STATUS_ACTIVE string = "active"
)

func (s *TasksService) Create(ctx context.Context, userID string, input types.TasksCreateDTO) error {
	stats, err := checkStatus(input.Status)
	if err != nil {
		return err
	}

	task := model.Task{
		UserID:   userID,
		Title:    input.Title,
		Status:   stats,
		Priority: input.Priority,
		Order:    input.Order,
	}
	err = s.repository.CreateTask(ctx, task)
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
	stats, err := checkStatus(input.Status)
	if err != nil {
		return types.TaskDTO{}, err
	}

	task, err := s.repository.UpdateById(ctx, userID, model.Task{
		ID:       input.ID,
		Title:    input.Title,
		Priority: input.Priority,
		Status:   stats,
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
		Status:   stats,
		Order:    task.Order,
	}, err
}

func (s *TasksService) ChangeStatus(ctx context.Context, userID, taskID, status string) error {
	stats, err := checkStatus(status)
	if err != nil {
		return err
	}

	err = s.repository.ChangeStatus(ctx, userID, taskID, stats)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func (s *TasksService) DeleteAll(ctx context.Context, userID string, status string) error {
	stats, err := checkStatus(status)
	if err != nil {
		return err
	}
	err = s.repository.DeleteAll(ctx, userID, stats)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func checkStatus(input string) (string, error) {
	var status string
	switch input {
	case STATUS_ACTIVE:
		status = STATUS_ACTIVE
	case STATUS_DELITE:
		status = STATUS_DELITE
	case STATUS_DONE:
		status = STATUS_DONE
	default:
		return status, errors.New("incorrect status")
	}
	return status, nil
}
