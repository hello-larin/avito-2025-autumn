package team

type team struct {
	TeamName string       `json:"team_name"`
	Members  []teamMember `json:"members"`
}

type teamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamCreatedResponse struct {
	Team team `json:"team"`
}
