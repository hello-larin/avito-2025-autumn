package team

type team struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []teamMember `json:"members"   validate:"required"`
}

type teamMember struct {
	UserID   string `json:"user_id"   validate:"required"`
	Username string `json:"username"  validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type teamCreatedResponse struct {
	Team team `json:"team"`
}
