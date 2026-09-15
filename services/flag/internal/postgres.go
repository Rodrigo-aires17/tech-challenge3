package internal

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository persiste as flags no RDS PostgreSQL (flag_db).
// As credenciais chegam via variáveis de ambiente (populadas a partir de um
// Kubernetes Secret sincronizado do AWS Secrets Manager), nunca em texto
// plano no repositório.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	repo := &PostgresRepository{pool: pool}
	if err := repo.migrate(ctx); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresRepository) migrate(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS feature_flags (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			enabled BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

// DatabaseURLFromEnv monta a connection string a partir das variáveis
// populadas pelo Secret sincronizado (DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, DB_PASSWORD).
func DatabaseURLFromEnv() string {
	return "postgres://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") +
		"?sslmode=require"
}

func (r *PostgresRepository) Create(ctx context.Context, name, description string, enabled bool) (Flag, error) {
	var flag Flag
	err := r.pool.QueryRow(ctx, `
		INSERT INTO feature_flags (name, description, enabled)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, enabled, created_at, updated_at`,
		name, description, enabled,
	).Scan(&flag.ID, &flag.Name, &flag.Description, &flag.Enabled, &flag.CreatedAt, &flag.UpdatedAt)
	return flag, err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Flag, error) {
	var flag Flag
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, description, enabled, created_at, updated_at
		FROM feature_flags WHERE id = $1`, id,
	).Scan(&flag.ID, &flag.Name, &flag.Description, &flag.Enabled, &flag.CreatedAt, &flag.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Flag{}, ErrNotFound
	}
	return flag, err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Flag, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, enabled, created_at, updated_at FROM feature_flags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []Flag
	for rows.Next() {
		var flag Flag
		if err := rows.Scan(&flag.ID, &flag.Name, &flag.Description, &flag.Enabled, &flag.CreatedAt, &flag.UpdatedAt); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	return flags, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, id string, enabled bool, description string) (Flag, error) {
	var flag Flag
	err := r.pool.QueryRow(ctx, `
		UPDATE feature_flags SET enabled = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, enabled, created_at, updated_at`,
		id, enabled, description,
	).Scan(&flag.ID, &flag.Name, &flag.Description, &flag.Enabled, &flag.CreatedAt, &flag.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Flag{}, ErrNotFound
	}
	return flag, err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM feature_flags WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
