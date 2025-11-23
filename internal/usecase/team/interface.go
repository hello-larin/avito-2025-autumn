package team

import (
	"context"

	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type teamRepository interface {
	GetTeamByName(ctx context.Context, name string) (*models.TeamDB, error)
	CreateTeam(ctx context.Context, name string) (*models.TeamDB, error)
}

type userRepository interface {
	AddUserToTeam(ctx context.Context, teamName string, user models.UserDB) (*models.UserDB, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]models.UserDB, error)
}
