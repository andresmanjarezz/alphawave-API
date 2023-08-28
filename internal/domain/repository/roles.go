package repository

import "context"

type RolesRepository interface {
	Create(ctx context.Context, teamID string)
}
