package models

import "time"

type TeamDB struct {
	TeamName string     `db:"team_name"`
	Created  *time.Time `db:"created_at"`
}
