package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type TasksRepository interface {
	Create(ctx context.Context, input model.Task) error
	GetById(ctx context.Context, userID, taskID string) (model.Task, error)
	GetAll(ctx context.Context, userID string) ([]model.Task, error)
}
