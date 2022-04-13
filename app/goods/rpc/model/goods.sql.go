// Code generated by sqlc. DO NOT EDIT.
// source: goods.sql

package model

import (
	"context"
	"database/sql"
	"time"
)

const createGoods = `-- name: CreateGoods :one
INSERT INTO "goods"(name, price)
VALUES ($1, $2) returning id, created_at, updated_at, deleted_at, name, price
`

type CreateGoodsParams struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (q *Queries) CreateGoods(ctx context.Context, arg CreateGoodsParams) (Good, error) {
	row := q.db.QueryRowContext(ctx, createGoods, arg.Name, arg.Price)
	var i Good
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Price,
	)
	return i, err
}

const deleteGoods = `-- name: DeleteGoods :one
UPDATE "goods"
set deleted_at =$1
where id = $2 returning id, created_at, updated_at, deleted_at, name, price
`

type DeleteGoodsParams struct {
	DeletedAt sql.NullTime `json:"deleted_at"`
	ID        int64        `json:"id"`
}

func (q *Queries) DeleteGoods(ctx context.Context, arg DeleteGoodsParams) (Good, error) {
	row := q.db.QueryRowContext(ctx, deleteGoods, arg.DeletedAt, arg.ID)
	var i Good
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Price,
	)
	return i, err
}

const getGoodsByID = `-- name: GetGoodsByID :one
SELECT id, created_at, updated_at, deleted_at, name, price
FROM "goods"
WHERE id = $1
  and deleted_at IS NULL
`

func (q *Queries) GetGoodsByID(ctx context.Context, id int64) (Good, error) {
	row := q.db.QueryRowContext(ctx, getGoodsByID, id)
	var i Good
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Price,
	)
	return i, err
}

const getGoodsByName = `-- name: GetGoodsByName :one
SELECT id, created_at, updated_at, deleted_at, name, price
FROM "goods"
WHERE name = $1
  and deleted_at IS NULL
`

func (q *Queries) GetGoodsByName(ctx context.Context, name string) (Good, error) {
	row := q.db.QueryRowContext(ctx, getGoodsByName, name)
	var i Good
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Price,
	)
	return i, err
}

const updateGoods = `-- name: UpdateGoods :one
UPDATE "goods"
SET updated_at = $1,
    name       = $2,
    price      = $3
WHERE id = $4
  and deleted_at IS NULL returning id, created_at, updated_at, deleted_at, name, price
`

type UpdateGoodsParams struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	ID        int64     `json:"id"`
}

func (q *Queries) UpdateGoods(ctx context.Context, arg UpdateGoodsParams) (Good, error) {
	row := q.db.QueryRowContext(ctx, updateGoods,
		arg.UpdatedAt,
		arg.Name,
		arg.Price,
		arg.ID,
	)
	var i Good
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Price,
	)
	return i, err
}