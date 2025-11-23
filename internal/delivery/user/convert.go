package user

import (
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

func toDTO(u *models.UserDB) user {
	return user{
		UserID:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func toSetActiveResponse(user *models.UserDB) userUpdateResponse {
	return userUpdateResponse{
		User: toDTO(user),
	}
}

func toPRDTO(pr *models.PullRequestShortDB) pullRequestShort {
	return pullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}

func toPRResponse(userID string, prs []models.PullRequestShortDB) userPullRequestsResponse {
	shorts := make([]pullRequestShort, len(prs))
	for i, pr := range prs {
		shorts[i] = toPRDTO(&pr)
	}
	return userPullRequestsResponse{
		UserID:       userID,
		PullRequests: shorts,
	}
}
