package user

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type Usecase struct {
	userRepo userRepository
	prRepo   prRepository
	manager  *manager.Manager
}

func New(userRepo userRepository, prRepo prRepository, manager *manager.Manager) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		prRepo:   prRepo,
		manager:  manager,
	}
}

func (uc *Usecase) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.UserDB, error) {
	_, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userDB, err := uc.userRepo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return userDB, nil
}

func (uc *Usecase) GetUserAssignedPRs(ctx context.Context, userID string) ([]models.PullRequestShortDB, error) {
	_, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	prs, err := uc.prRepo.GetUserPRs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
