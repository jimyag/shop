package model

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SqlStore) ExecTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}
