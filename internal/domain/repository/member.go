package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type MemberRepository interface {
	GetMembers(ctx context.Context, teamID string) ([]model.Member, error)
}
