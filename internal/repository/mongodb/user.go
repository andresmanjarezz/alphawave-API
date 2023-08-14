package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection(usersCollection),
	}
}

func (r *UserRepository) Create(ctx context.Context, input model.User) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := r.db.InsertOne(nCtx, input); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByСredentials(ctx context.Context, email, password string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	filter := bson.M{"email": email, "password": password}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}
		return user, err
	}

	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	filter := bson.M{"email": email}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}
		return user, err
	}

	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) GetUserById(ctx context.Context, userID string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}
	filter := bson.M{"_id": ObjectID}

	res := r.db.FindOne(nCtx, filter)

	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}

		return user, err

	}
	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) ChangeVerificationCode(ctx context.Context, email string, input model.UserVerificationPayload) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.db.UpdateOne(nCtx, bson.M{"email": email}, bson.M{"$set": bson.M{"verification.verificationCode": input.VerificationCode, "verification.verificationCodeExpiresTime": input.VerificationCodeExpiresTime}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepository) SetSession(ctx context.Context, userID string, session model.Session, lastVisitTime time.Time) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$set": bson.M{"session": session, "lastVisitTime": lastVisitTime}})

	return err
}

func (r *UserRepository) Verify(ctx context.Context, verificationCode string) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(nCtx, bson.M{"verification.verificationCode": verificationCode}, bson.M{"$set": bson.M{"verification.verified": true, "verification.verificationCode": ""}})
	return err
}

func (r *UserRepository) ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	res, err := r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID, "password": oldPassword}, bson.M{"$set": bson.M{"password": newPassword}})

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) SetForgotPassword(ctx context.Context, email string, input model.ForgotPasswordPayload) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(nCtx, bson.M{"email": email}, bson.M{"$set": bson.M{"forgotPasswordToken": input}})

	return err
}

func (r *UserRepository) GetUserByVerificationCode(ctx context.Context, hash string) (model.User, error) {
	var user model.User
	filter := bson.M{"verification.verificationCode": hash}

	res := r.db.FindOne(ctx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}

		return user, err

	}
	if err := res.Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var user model.User

	if err := r.db.FindOne(nCtx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresTime":  bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, apperrors.ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) IsDuplicate(ctx context.Context, email string) (bool, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	filter := bson.M{"email": email}

	count, err := r.db.CountDocuments(nCtx, filter)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil

}

func (r *UserRepository) GetByForgotPasswordToken(ctx context.Context, token, tokenResult string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User

	filter := bson.M{"forgotPasswordToken.token": token, "forgotPasswordToken.resultToken": tokenResult, "forgotPasswordToken.tokenExpiresTime": bson.M{"$gt": time.Now()}}

	res := r.db.FindOne(nCtx, filter)

	if err := res.Decode(&user); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, apperrors.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) ResetPassword(ctx context.Context, token, email, password string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"forgotPasswordToken.token": token, "email": email}

	res, err := r.db.UpdateOne(nCtx, filter, bson.M{"$set": bson.M{"password": password}, "$unset": bson.M{"forgotPasswordToken": ""}})

	if res.MatchedCount == 0 {
		return apperrors.ErrUserNotFound
	}
	return err
}
