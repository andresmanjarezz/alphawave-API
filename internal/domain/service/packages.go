package service

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
)

type PackagesService struct {
	repository repository.PackagesRepository
}

func NewPackagesService(repository repository.PackagesRepository) *PackagesService {
	return &PackagesService{
		repository: repository,
	}
}

// func (s *PackagesService) CreateDefaultPackages() error {
// 	var feature model.Feature

// 	defaultPackage := model.Package{
// 		Name:        "Start",
// 		Description: "Package for start",
// 		Price:       20,
// 		Currency:    "usd",
// 	}
// }
