package user

type userSetActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type user struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type userUpdateResponse struct {
	User user `json:"user"`
}

type pullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type userPullRequestsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []pullRequestShort `json:"pull_requests"`
}
