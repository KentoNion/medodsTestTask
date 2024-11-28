package Store

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type Store struct {
	conn *sql.DB
}

func (s Store) Save(ctx context.Context, token string, userID string) error {
	query := "INSERT INTO tokens (user_id, token) VALUES (?, ?) ON CONFLICT DO UPDATE token = ?"
	_, err := s.conn.ExecContext(ctx, query, userID, token, token)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, token string) (bool, error) {
	query := "SELECT 1 FROM tokens where token = ?"
	rows, err := s.conn.QueryContext(ctx, query, token)
	if err != nil {
		return false, errors.Wrap(err, "failed to query")
	}
	return rows.Next(), nil
}
