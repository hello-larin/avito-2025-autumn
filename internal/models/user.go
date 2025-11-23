package models

type UserDB struct {
	UserID   string `db:"id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}
