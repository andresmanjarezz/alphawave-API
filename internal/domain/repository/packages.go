package repository

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type PackagesRepository interface {
	// GetByID(ctx context.Context, id string)
	CreateDefaultPackages(packages []model.Package) error
}
