package mongodb

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackagesRepository struct {
	db *mongo.Collection
}

func NewPackagesRepository(db *mongo.Database) *PackagesRepository {
	return &PackagesRepository{
		db: db.Collection(packagesCollection),
	}
}

func (r *PackagesRepository) CreateDefaultPackages(packages []model.Package) error {
	nCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := r.db.InsertOne(nCtx, packages); err != nil {
		return err
	}
	return nil
}
