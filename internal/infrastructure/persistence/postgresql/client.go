package postgresql

import (
	"context"
	"database/sql"

	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/postgresql/sqlc"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type client struct {
	db *sql.DB
}

func NewClient(ctx context.Context, database_url string) (*client, error) {
	db, err := sql.Open("pgx", database_url)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &client{
		db: db,
	}, nil
}

func (c *client) Queries() *sqlc.Queries {
	return sqlc.New(c.db)
}
