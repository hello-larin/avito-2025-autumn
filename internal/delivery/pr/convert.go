package pr

import (
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

func (r *prCreateRequest) ToModel() models.PullRequestDB {
	return models.PullRequestDB{
		PullRequestID:   r.PullRequestID,
		PullRequestName: r.PullRequestName,
		AuthorID:        r.AuthorID,
	}
}

func toDTO(pr *models.PullRequestDB, assignedReviewers []string) pullRequest {
	return pullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         nil,
		MergedAt:          pr.MergedAt,
	}
}

func toShortDTO(pr *models.PullRequestShortDB, assignedReviewers []string) pullRequest {
	return pullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: assignedReviewers,
	}
}

func toCreateResponse(pr *models.PullRequestShortDB, assignedReviewers []string) prCreateResponse {
	return prCreateResponse{
		PullRequest: toShortDTO(pr, assignedReviewers),
	}
}

func toMergeResponse(pr *models.PullRequestDB, assignedReviewers []string) prMergeResponse {
	return prMergeResponse{
		PullRequest: toDTO(pr, assignedReviewers),
	}
}

func toReassignResponse(pr *models.PullRequestDB, assignedReviewers []string, newReviewerID string) prReassignResponse {
	return prReassignResponse{
		PullRequest: toDTO(pr, assignedReviewers),
		ReplacedBy:  newReviewerID,
	}
}
