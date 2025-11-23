package pr

import (
	"context"

	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type prRepository interface {
	GetPRByID(ctx context.Context, prID string) (*models.PullRequestDB, error)
	CreatePR(ctx context.Context, pr *models.PullRequestDB) (*models.PullRequestShortDB, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequestDB, error)
	GetAllPRReviewers(ctx context.Context, prID string) ([]string, error)
	GetActivePRReviewers(ctx context.Context, prID string) ([]string, error)
	AssignReviewer(ctx context.Context, prID, userID string) error
	UnassignReviewer(ctx context.Context, prID, userID string) error
}

type userRepository interface {
	GetUserByID(ctx context.Context, userID string) (*models.UserDB, error)
	GetActiveTeamMembers(ctx context.Context, teamName string, limit int) ([]models.UserDB, error)
}
