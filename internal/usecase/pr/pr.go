package pr

import (
	"context"
	"errors"
	"slices"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"
)

type Usecase struct {
	prRepo   prRepository
	userRepo userRepository
	manager  *manager.Manager
}

func New(prRepo prRepository, userRepo userRepository, manager *manager.Manager) *Usecase {
	return &Usecase{
		prRepo:   prRepo,
		userRepo: userRepo,
		manager:  manager,
	}
}

func (uc *Usecase) CreatePullRequest(
	ctx context.Context,
	pr *models.PullRequestDB,
) (*models.PullRequestShortDB, []string, error) {
	author, err := uc.userRepo.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, nil, err
	}
	_, err = uc.prRepo.GetPRByID(ctx, pr.PullRequestID)
	if err == nil {
		return nil, nil, customerror.ErrPRExists
	}
	if !errors.Is(err, customerror.ErrNotFound) {
		return nil, nil, err
	}
	var createdPR *models.PullRequestShortDB
	var insertedReviewers []string
	err = uc.manager.Do(ctx, func(ctx context.Context) error {

		createdPR, err = uc.prRepo.CreatePR(ctx, pr)
		if err != nil {
			return err
		}

		// Может выпасть ещё автор поэтому +1
		activeMembers, err := uc.userRepo.GetActiveTeamMembers(ctx, author.TeamName, 3)
		if err != nil {
			return err
		}

		// По заданию до 2 ревьюверов на ПР
		for _, user := range activeMembers {
			if len(insertedReviewers) == 2 {
				break
			}
			if user.UserID == pr.AuthorID {
				continue
			}
			err = uc.prRepo.AssignReviewer(ctx, createdPR.PullRequestID, user.UserID)
			insertedReviewers = append(insertedReviewers, user.UserID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return createdPR, insertedReviewers, nil
}

func (uc *Usecase) MergePullRequest(
	ctx context.Context,
	prID string,
) (*models.PullRequestDB, []string, error) {
	_, err := uc.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, nil, err
	}
	mergedPR, err := uc.prRepo.MergePR(ctx, prID)
	if err != nil {
		return nil, nil, err
	}

	reviewers, err := uc.prRepo.GetActivePRReviewers(ctx, prID)
	if err != nil {
		return nil, nil, err
	}

	return mergedPR, reviewers, nil
}

func (uc *Usecase) ReassignReviewer(
	ctx context.Context,
	prID, oldReviewerID string,
) (*models.PullRequestDB, []string, string, error) {
	pr, err := uc.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, nil, "", customerror.ErrNotFound
	}

	if pr.Status == "MERGED" {
		return nil, nil, "", customerror.ErrPRMerged
	}

	activeReviewers, err := uc.prRepo.GetActivePRReviewers(ctx, prID)
	if err != nil {
		return nil, nil, "", err
	}

	isAssigned := slices.Contains(activeReviewers, oldReviewerID)
	if !isAssigned {
		return nil, nil, "", customerror.ErrNotAssigned
	}

	author, err := uc.userRepo.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, nil, "", customerror.ErrNotFound
	}

	// Тут можно сделать лучше, но пока так
	activeMembers, err := uc.userRepo.GetActiveTeamMembers(ctx, author.TeamName, 10)
	if err != nil {
		return nil, nil, "", err
	}
	if len(activeMembers) == 0 {
		return nil, nil, "", customerror.ErrNoCandidate
	}

	allReviewers, err := uc.prRepo.GetAllPRReviewers(ctx, prID)
	if err != nil {
		return nil, nil, "", err
	}

	// Мапа для поиска кандидата. true - пропускаем. изначальное значение bool = false
	exclude := make(map[string]bool)
	exclude[pr.AuthorID] = true
	for _, r := range allReviewers {
		exclude[r] = true
	}
	var newReviewerID string
	for _, member := range activeMembers {
		if !exclude[member.UserID] {
			newReviewerID = member.UserID
			break
		}
	}
	if newReviewerID == "" {
		return nil, nil, "", customerror.ErrNoCandidate
	}

	err = uc.manager.Do(ctx, func(ctx context.Context) error {

		if err = uc.prRepo.UnassignReviewer(ctx, prID, oldReviewerID); err != nil {
			return err
		}

		if err = uc.prRepo.AssignReviewer(ctx, prID, newReviewerID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, "", err
	}

	updatedReviewers, err := uc.prRepo.GetActivePRReviewers(ctx, prID)
	if err != nil {
		return nil, nil, "", err
	}

	return pr, updatedReviewers, newReviewerID, nil
}
