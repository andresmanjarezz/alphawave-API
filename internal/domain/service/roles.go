package service

import "github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"

type RolesService struct {
	repository repository.RolesRepository
}

func NewRolesService(repository repository.RolesRepository) *RolesService {
	return &RolesService{
		repository: repository,
	}
}
