package mongodb

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamsRepository struct {
	db *mongo.Collection
}

func NewTeamsRepository(db *mongo.Database) *TeamsRepository {
	return &TeamsRepository{
		db: db.Collection(teamsCollection),
	}
}

func (r *TeamsRepository) CreateTeam(ctx context.Context, input model.Team) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := r.db.InsertOne(nCtx, input); err != nil {
		return err
	}

	return nil
}
