package pr

import "time"

type pullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
	CreatedAt         *time.Time `json:"-"`
	//  MergedAt используетя, а CreatedAt
	// вообще нигде не используется
}

type prCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
}

type prCreateResponse struct {
	PullRequest pullRequest `json:"pr"`
}

type prMergeRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

type prMergeResponse struct {
	PullRequest pullRequest `json:"pr"`
}

type prReassignRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldUserID     string `json:"old_reviewer_id" validate:"required"`
}

type prReassignResponse struct {
	PullRequest pullRequest `json:"pr"`
	ReplacedBy  string      `json:"replaced_by,omitempty"`
}
