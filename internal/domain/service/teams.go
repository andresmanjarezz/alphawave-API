package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type TeamsService struct {
	repository       repository.TeamsRepository
	userRepository   repository.UserRepository
	memberRepository repository.MemberRepository
	rolesService     RolesService
}

func NewTeamsService(repository repository.TeamsRepository, userRepository repository.UserRepository, memberRepository repository.MemberRepository, rolesService RolesService) *TeamsService {
	return &TeamsService{
		repository:       repository,
		userRepository:   userRepository,
		memberRepository: memberRepository,
		rolesService:     rolesService,
	}
}

func (s *TeamsService) Create(ctx context.Context, userID string, input types.CreateTeamsDTO) error {

	team := model.Team{
		TeamName: input.TeamName,
		JobTitle: input.JobTitle,
		OwnerID:  userID,
	}
	id, err := s.repository.CreateTeam(ctx, team)
	if err != nil {
		return err
	}
	err = s.rolesService.Create(ctx, id)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetUserById(ctx, userID)

	if err != nil {
		return err
	}
	roles := make([]string, 0, 1)

	roles = append(roles, model.ROLE_OWNER)
	if err := s.memberRepository.CreateMember(ctx, id, model.Member{
		TeamID: id,
		UserID: user.ID,
		Email:  user.Email,
		Status: USER_STATUS_ACTIVE,
		Roles:  roles,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	return nil
}

func (s *TeamsService) UpdateTeamSettings(ctx context.Context, teamID string, input types.UpdateTeamSettingsDTO) error {

	if err := s.repository.UpdateTeamSettings(ctx, teamID, model.UpdateTeamSettingsInput{
		LogoURL:               input.LogoURL,
		UserActivityIndicator: input.UserActivityIndicator,
		DisplayLinkPreview:    input.DisplayLinkPreview,
		DisplayFilePreview:    input.DisplayFilePreview,
		EnableGifs:            input.EnableGifs,
		ShowWeekends:          input.ShowWeekends,
		FirstDayOfWeek:        input.FirstDayOfWeek,
	}); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}

	return nil
}

func (s *TeamsService) GetTeamByID(ctx context.Context, teamID string) (model.Team, error) {
	return s.repository.GetTeamByID(ctx, teamID)
}

func (s *TeamsService) GetTeamsByUser(ctx context.Context, userID string) ([]model.Team, error) {
	members, err := s.memberRepository.GetMembersByUserID(ctx, userID)

	if err != nil {
		return []model.Team{}, err
	}

	teamsIds := make([]string, 0, len(members))
	fmt.Printf("members: %v", members)
	for _, member := range members {
		fmt.Printf("teamID: %s", member.TeamID)
		teamsIds = append(teamsIds, member.TeamID)
	}

	teams, err := s.repository.GetTeamsByIds(ctx, teamsIds)

	if err != nil {
		return []model.Team{}, err
	}

	return teams, nil
}

func (s *TeamsService) GetTeamsByIds(ctx context.Context, ids []string) ([]model.Team, error) {
	teams, err := s.repository.GetTeamsByIds(ctx, ids)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []model.Team{}, apperrors.ErrDocumentNotFound
		}
		return []model.Team{}, err
	}
	return teams, nil
}
