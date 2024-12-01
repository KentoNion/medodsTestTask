package store

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Store struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
	sm sqluct.Mapper
}

func NewDB(db *sqlx.DB) *Store {
	return &Store{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		sm: sqluct.Mapper{Dialect: sqluct.DialectPostgres},
	}
}

func (s Store) Save(ctx context.Context, token string, userID string, ip string) error {
	query := "INSERT INTO tokens (user_id, token, ip) VALUES ($1, $2, $3) ON CONFLICT (token) DO UPDATE SET token = $2, ip = $3"
	_, err := s.db.ExecContext(ctx, query, userID, token, ip)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, token string) (bool, error) {
	query := "SELECT 1 FROM tokens where token = $1"
	rows, err := s.db.QueryContext(ctx, query, token)
	defer rows.Close()
	if err != nil {
		return false, errors.Wrap(err, "failed to query")
	}
	return rows.Next(), nil
}

func (s Store) Delete(ctx context.Context, token string) error {
	query := "DELETE FROM tokens WHERE token = $1"
	rows, err := s.db.ExecContext(ctx, query, token)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	if affected, _ := rows.RowsAffected(); affected == 0 {
		return errors.New("failed to delete")
	}
	return nil
}
