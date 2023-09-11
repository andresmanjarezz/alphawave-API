package service

import "context"

type SubscriptionService struct {
	UserService  UserServiceI
	TeamsService TeamsServiceI
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}

func (s *SubscriptionService) Create(ctx context.Context, userID string, packageID string, teamID string) (string, error) {
	// todo

	return "", nil
}
