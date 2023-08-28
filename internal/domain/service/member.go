package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

const (
	USER_STATUS_ACTIVE   = "ACTIVE"
	USER_STATUS_INACTIVE = "INACTIVE"
	USER_STATUS_PENDING  = "PENDING"
)

type MemberService struct {
	repository     repository.MemberRepository
	userRepository repository.UserRepository
}

func NewMemberService(repository repository.MemberRepository, userRepository repository.UserRepository) *MemberService {
	return &MemberService{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *MemberService) GetMembers(ctx context.Context, teamID string, query types.GetUsersByQuery) ([]types.MemberDTO, error) {
	members, err := s.repository.GetMembers(ctx, teamID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.MemberDTO{}, apperrors.ErrDocumentNotFound
		}
		return []types.MemberDTO{}, err
	}
	var userIds = make([]string, len(members))

	for _, member := range members {
		userIds = append(userIds, member.UserID)
	}

	users, err := s.userRepository.GetUsersByQuery(ctx, userIds, model.GetUsersByQuery{PaginationQuery: model.PaginationQuery(query.PaginationQuery)})
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.MemberDTO{}, apperrors.ErrDocumentNotFound
		}
		return []types.MemberDTO{}, err
	}

	var membersOutput = make([]types.MemberDTO, len(members))

	for i := range members {
		membersOutput = append(membersOutput, types.MemberDTO{
			MemberID:      users[i].ID,
			FirstName:     users[i].FirstName,
			LastName:      users[i].LastName,
			Email:         users[i].Email,
			LastVisitTime: users[i].LastVisitTime,
			Roles:         members[i].Roles,
		})
	}
	return membersOutput, nil
}

// func (s *MemberService) UserInvite(ctx context.Context, teamID string, email string, role string) error {

// 	user, err := s.userRepository.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		if errors.Is(err, apperrors.ErrUserNotFound) {

// 		}
// 		return err
// 	}

// 	err := s.repository.CreateMember(ctx, teamID)
// 	if err != nil {
// 		if errors.Is(err, apperrors.ErrDocumentNotFound) {
// 			return apperrors.ErrDocumentNotFound
// 		}
// 		return err
// 	}
// }
