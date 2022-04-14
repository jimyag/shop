-- name: CreateInventory :one
INSERT INTO "inventory"(goods_id,
                        sticks,
                        version)
VALUES ($1, $2, $3)
returning *;

-- name: GetInventoryByGoodsID :one
SELECT *
FROM "inventory"
WHERE goods_id = $1
LIMIT 1;


-- name: UpdateInventory :one
update "inventory"
set updated_at = $1,
    sticks     = sticks + sqlc.arg(counts)
where goods_id = $2
returning *;

