package core_postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := NewConfigMust()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_DB,
	)
	pool, err := pgxpool.New(
		ctx,
		connStr,
	)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func ClosePool(p *pgxpool.Pool) {
	p.Close()
}
