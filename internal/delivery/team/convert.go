package team

import "github.com/hello-larin/avito-2025-autumn/internal/models"

func (m *teamMember) toModel() models.UserDB {
	return models.UserDB{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

func toCreateResponse(teamName *models.TeamDB, members []models.UserDB) teamCreatedResponse {
	return teamCreatedResponse{
		Team: toDTO(teamName, members),
	}
}

func toDTO(teamName *models.TeamDB, members []models.UserDB) team {
	result := team{
		TeamName: teamName.TeamName,
	}
	for _, user := range members {
		result.Members = append(result.Members, toUserDTO(user))
	}
	return result
}

func toUserDTO(member models.UserDB) teamMember {
	return teamMember{
		UserID:   member.UserID,
		Username: member.Username,
		IsActive: member.IsActive,
	}
}
