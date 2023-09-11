package repository

import "context"

type PackagesRepository interface {
	GetByID(ctx context.Context, id string)
}
