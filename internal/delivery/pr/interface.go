package pr

import (
	"github.com/hello-larin/avito-2025-autumn

	"avito-2025-autumn/internal/models"
)

type usecase interface {
	CreatePullRequest(ctx context.Context, pr *models.PullRequestDB) (*models.PullRequestShortDB, []string, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequestDB, []string, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*models.PullRequestDB, []string, string, error)
}
