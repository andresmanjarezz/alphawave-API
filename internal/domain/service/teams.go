package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type TeamsService struct {
	repository     repository.TeamsRepository
	userRepository repository.UserRepository
}

func NewTeamsService(repository repository.TeamsRepository, userRepository repository.UserRepository) *TeamsService {
	return &TeamsService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *TeamsService) Create(ctx context.Context, userEmail string, input types.CreateTeamsDTO) error {
	user, err := s.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
	}

	team := model.Team{
		TeamName: input.TeamName,
		JobTitle: input.JobTitle,
		OwnerID:  user.ID,
	}
	if err := s.repository.CreateTeam(ctx, team); err != nil {
		return err
	}
	return nil
}

func (s *TeamsService) GetAllByIds(ctx context.Context, ids []string) (types.TeamsDTO, error) {
	return types.TeamsDTO{}, nil
}
