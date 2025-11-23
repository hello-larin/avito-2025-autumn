package models

import "time"

type PullRequestDB struct {
	PullRequestID   string     `db:"id"`
	PullRequestName string     `db:"name"`
	AuthorID        string     `db:"author_id"`
	Status          string     `db:"status"`
	CreatedAt       *time.Time `db:"created_at"`
	MergedAt        *time.Time `db:"merged_at"`
}

type PullRequestShortDB struct {
	PullRequestID   string `db:"id"`
	PullRequestName string `db:"name"`
	AuthorID        string `db:"author_id"`
	Status          string `db:"status"`
}
