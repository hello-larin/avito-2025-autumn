package team

import (
	"context"
	"errors"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type Usecase struct {
	teamRepo teamRepository
	userRepo userRepository
}

func New(teamRepo teamRepository, userRepo userRepository) *Usecase {
	return &Usecase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (uc *Usecase) CreateTeam(
	ctx context.Context,
	teamName string,
	members []models.UserDB,
) (*models.TeamDB, []models.UserDB, error) {
	_, err := uc.teamRepo.GetTeamByName(ctx, teamName)
	if err == nil {
		return nil, nil, customerror.ErrTeamExists
	}
	if !errors.Is(err, customerror.ErrNotFound) {
		return nil, nil, err
	}
	team, err := uc.teamRepo.CreateTeam(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}
	var createdUsers []models.UserDB
	for _, member := range members {
		var createdUser *models.UserDB
		createdUser, err = uc.userRepo.AddUserToTeam(ctx, teamName, member)
		if err != nil {
			return nil, nil, err
		}
		createdUsers = append(createdUsers, *createdUser)
	}
	return team, createdUsers, nil
}

func (uc *Usecase) GetTeam(ctx context.Context, teamName string) (*models.TeamDB, []models.UserDB, error) {
	team, err := uc.teamRepo.GetTeamByName(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}

	members, err := uc.userRepo.GetTeamMembers(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}

	return team, members, nil
}
