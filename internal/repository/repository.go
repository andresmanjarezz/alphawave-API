package repository

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	User  repository.UserRepository
	Tasks repository.TasksRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User:  mongodb.NewUserRepository(db),
		Tasks: mongodb.NewTasksRepository(db),
	}
}
