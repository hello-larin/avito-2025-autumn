package team

import (
	"context"

	customerror "github.com/hello-larin/avito-2025-autumn/internal/error"
	"github.com/hello-larin/avito-2025-autumn/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type Repository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func New(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *Repository {
	return &Repository{
		db:     db,
		getter: getter,
	}
}

func (r *Repository) GetTeamByName(ctx context.Context, name string) (*models.TeamDB, error) {
	const query = `
		SELECT team_name, created_at 
		FROM teams 
		WHERE team_name = $1
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	row, err := conn.Query(ctx, query, name)
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
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	row, err := conn.Query(ctx, query, name)
	if err != nil {
		return nil, err
	}
	team, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.TeamDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}

	return &team, nil
}
