package user

import (
	"github.com/hello-larin/avito-2025-autumn

	"avito-2025-autumn/internal/models"
)

type userRepository interface {
	GetUserByID(ctx context.Context, userID string) (*models.UserDB, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.UserDB, error)
}

type prRepository interface {
	GetUserPRs(ctx context.Context, userID string) ([]models.PullRequestShortDB, error)
}
