package Store

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

func (s Store) Save(ctx context.Context, token string, userID string) error {
	query := "INSERT INTO tokens (user_id, token) VALUES (?, ?) ON CONFLICT DO UPDATE token = ?"
	_, err := s.db.ExecContext(ctx, query, userID, token, token)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, token string) (bool, error) {
	query := "SELECT 1 FROM tokens where token = ?"
	rows, err := s.db.QueryContext(ctx, query, token)
	if err != nil {
		return false, errors.Wrap(err, "failed to query")
	}
	return rows.Next(), nil
}
