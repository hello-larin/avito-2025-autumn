package pr

import (
	"context"
	"log/slog"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) GetPRByID(ctx context.Context, id string) (*models.PullRequestDB, error) {
	const query = `
		SELECT id, name, author_id, status, created_at, merged_at
		FROM pull_requests 
		WHERE id = $1
	`
	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	pr, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.PullRequestDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	slog.Info("DB_PR", "QUERY", "GET_PR_BY_ID", "DATA", pr)
	return &pr, nil
}

func (r *PullRequestRepository) CreatePR(
	ctx context.Context,
	pr *models.PullRequestDB,
) (*models.PullRequestShortDB, error) {
	const query = `
	INSERT INTO pull_requests (id, name, author_id) 
	VALUES ($1, $2, $3)
	RETURNING id, name, author_id, status
	`
	row, err := r.db.Query(ctx, query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID)
	if err != nil {
		return nil, customerror.ErrPRExists
	}
	createdPR, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.PullRequestShortDB])
	if err != nil {
		return nil, err
	}
	slog.Info("DB_PR", "QUERY", "CREATE_PR", "DATA", createdPR)
	return &createdPR, err
}

func (r *PullRequestRepository) AssignReviewer(ctx context.Context, prID, userID string) error {
	const query = `
	INSERT INTO pr_reviewers (pr_id, user_id)
	VALUES ($1, $2)
	`
	_, err := r.db.Exec(ctx, query, prID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepository) UnassignReviewer(ctx context.Context, prID, userID string) error {
	const query = `
	UPDATE pr_reviewers 
	SET is_removed = TRUE 
	WHERE pr_id = $1 AND user_id = $2 AND is_removed = FALSE
	`
	row, err := r.db.Exec(ctx, query, prID, userID)
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return customerror.ErrNotAssigned
	}

	return nil
}

func (r *PullRequestRepository) MergePR(ctx context.Context, id string) (*models.PullRequestDB, error) {
	const query = `
	UPDATE pull_requests 
	SET status = 'MERGED' WHERE id = $1
	RETURNING id, name, author_id, status, created_at, merged_at
	`
	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	updatedPR, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.PullRequestDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	slog.Info("DB_PR", "QUERY", "MERGE_PR", "DATA", updatedPR)

	return &updatedPR, err
}

func (r *PullRequestRepository) GetActivePRReviewers(ctx context.Context, prID string) ([]string, error) {
	const query = `
		SELECT user_id 
		FROM pr_reviewers 
		WHERE pr_id = $1 AND is_removed = FALSE
	`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	reviewers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err = row.Scan(&n)
		return n, err
	})
	if err = rows.Err(); err != nil {
		return nil, err
	}
	slog.Info("DB_PR", "QUERY", "GET_ACTIVE_PR_REVIEWERS", "DATA", reviewers)

	return reviewers, nil
}

func (r *PullRequestRepository) GetAllPRReviewers(ctx context.Context, prID string) ([]string, error) {
	const query = `
		SELECT user_id 
		FROM pr_reviewers 
		WHERE pr_id = $1
	`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	reviewers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err = row.Scan(&n)
		return n, err
	})
	slog.Info("DB_PR", "QUERY", "GET_ALL_PR_REVIEWERS", "DATA", reviewers)

	return reviewers, nil
}

func (r *PullRequestRepository) GetUserPRs(ctx context.Context, userID string) ([]models.PullRequestShortDB, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.user_id = $1 AND prr.is_removed = FALSE
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	prs, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PullRequestShortDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	slog.Info("DB_PR", "QUERY", "GET_USER_PR", "DATA", prs)

	return prs, nil
}
