package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type client struct {
	db *pgx.Conn
}

func NewClient(ctx context.Context, database_url string) (*client, error) {
	connection, err := pgx.Connect(ctx, database_url)

	if err != nil {
		return nil, err
	}

	return &client{
		db: connection,
	}, nil
}
