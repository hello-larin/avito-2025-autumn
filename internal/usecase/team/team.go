package team

import (
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type Usecase struct {
	teamRepo teamRepository
	userRepo userRepository
	manager  *manager.Manager
}

func New(teamRepo teamRepository, userRepo userRepository, manager *manager.Manager) *Usecase {
	return &Usecase{
		teamRepo: teamRepo,
		userRepo: userRepo,
		manager:  manager,
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
	var createdUsers []models.UserDB
	var team *models.TeamDB
	err = uc.manager.Do(ctx, func(ctx context.Context) error {
		team, err = uc.teamRepo.CreateTeam(ctx, teamName)
		if err != nil {
			return err
		}
		for _, member := range members {
			var createdUser *models.UserDB
			createdUser, err = uc.userRepo.AddUserToTeam(ctx, teamName, member)
			if err != nil {
				return err
			}
			createdUsers = append(createdUsers, *createdUser)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
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
