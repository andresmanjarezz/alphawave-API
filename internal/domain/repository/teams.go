package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type TeamsRepository interface {
	CreateTeam(ctx context.Context, input model.Team) error
	// GetTeamByID(ctx context.Context, teamID string) (model.Team, error)
}
