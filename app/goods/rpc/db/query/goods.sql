-- name: GetGoodsByID :one
SELECT *
FROM "goods"
WHERE id = $1
  and deleted_at IS NULL
;

-- name: GetGoodsByName :one
SELECT *
FROM "goods"
WHERE name = $1
  and deleted_at IS NULL
;

-- name: CreateGoods :one
INSERT INTO "goods"(name, price)
VALUES ($1, $2) returning *;

-- name: DeleteGoods :one
UPDATE "goods"
set deleted_at =$1
where id = $2 returning *;

-- name: UpdateGoods :one
UPDATE "goods"
SET updated_at = $1,
    name       = $2,
    price      = $3
WHERE id = $4
  and deleted_at IS NULL returning *;
