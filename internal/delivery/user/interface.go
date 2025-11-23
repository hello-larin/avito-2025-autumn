package user

import (
	"context"

	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type usecase interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.UserDB, error)
	GetUserAssignedPRs(ctx context.Context, userID string) ([]models.PullRequestShortDB, error)
}
