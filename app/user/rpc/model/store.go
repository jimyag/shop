package model

import (
	"context"
	"database/sql"
	"fmt"
)

//
// Store
//  @Description: 存储的接口，可以实现这个接口将数据存储在不同的介质中
//
type Store interface {
	ExecTx(ctx context.Context, fn func(queries *Queries) error) error
	Querier
}

//
// SqlStore
//  @Description: 使用数据库进行存储
//
type SqlStore struct {
	*Queries
	db *sql.DB
}

//
// NewSQLStore
//  @Description: 数据库存储的对象
//  @param db
//  @return Store
//
func NewSQLStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

//
// ExecTx
//  @Description: 执行事务
//  @receiver store
//  @param ctx
//  @param fn
//  @return error
//
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
