package user

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

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*models.UserDB, error) {
	const query = `
		SELECT id, username, team_name, is_active 
		FROM users 
		WHERE id = $1
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	row, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	return &user, nil
}

func (r *Repository) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.UserDB, error) {
	const query = `
		UPDATE users 
		SET is_active = $1 
		WHERE id = $2
		RETURNING id, username, team_name, is_active
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	row, err := conn.Query(ctx, query, isActive, userID)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	return &user, nil
}

func (r *Repository) GetTeamMembers(ctx context.Context, teamName string) ([]models.UserDB, error) {
	query := `
		SELECT id, username, team_name, is_active
		FROM users 
		WHERE team_name = $1
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	rows, err := conn.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserDB])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetActiveTeamMembers(
	ctx context.Context,
	teamName string,
	limit int,
) ([]models.UserDB, error) {
	query := `
		SELECT id, username, team_name, is_active
		FROM users 
		WHERE team_name = $1 AND is_active = TRUE
		ORDER BY RANDOM() 
		LIMIT $2
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	rows, err := conn.Query(ctx, query, teamName, limit)
	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserDB])
	if err != nil {
		return nil, customerror.ErrNotFound
	}
	return users, nil
}

func (r *Repository) AddUserToTeam(
	ctx context.Context,
	teamName string,
	user models.UserDB,
) (*models.UserDB, error) {
	const query = `
	INSERT INTO users (id, username, team_name, is_active) 
	 VALUES ($1, $2, $3, $4)
	ON CONFLICT (id)
	 DO UPDATE SET team_name = $3 
	RETURNING id, username, team_name, is_active
	`
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	rows, err := conn.Query(ctx, query, user.UserID, user.Username, teamName, user.IsActive)
	if err != nil {
		return nil, err
	}
	newUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserDB])
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}
