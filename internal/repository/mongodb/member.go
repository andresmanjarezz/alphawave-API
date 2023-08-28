package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MemberRepository struct {
	db *mongo.Collection
}

func NewMemberRepository(db *mongo.Database) *MemberRepository {
	return &MemberRepository{
		db: db.Collection(memberCollection),
	}
}

func (r *MemberRepository) GetMembers(ctx context.Context, teamID string) ([]model.Member, error) {
	nCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var members []model.Member

	cur, err := r.db.Find(nCtx, bson.M{"teamID": teamID})

	if err != nil {
		return []model.Member{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Member{}, apperrors.ErrDocumentNotFound
		}
		return []model.Member{}, err
	}

	if err := cur.All(nCtx, &members); err != nil {
		return []model.Member{}, err
	}

	return members, nil
}
