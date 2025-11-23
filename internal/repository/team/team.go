package team

import (
	"context"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTeamByName(ctx context.Context, name string) (*models.TeamDB, error) {
	const query = `
		SELECT team_name, created_at 
		FROM teams 
		WHERE team_name = $1
	`
	row, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, err
	}
	team, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.TeamDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	return &team, nil
}

func (r *Repository) CreateTeam(ctx context.Context, name string) (*models.TeamDB, error) {
	const query = `
	 INSERT INTO teams (team_name) 
	 VALUES ($1)
	 RETURNING team_name, created_at
	`
	row, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, err
	}
	team, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.TeamDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}

	return &team, nil
}
