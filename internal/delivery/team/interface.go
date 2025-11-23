package team

import (
	"context"

	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type usecase interface {
	CreateTeam(ctx context.Context, teamName string, members []models.UserDB) (*models.TeamDB, []models.UserDB, error)
	GetTeam(ctx context.Context, teamName string) (*models.TeamDB, []models.UserDB, error)
}
