-- name: GetCartListByUid :many
SELECT *
FROM "shopping_cart"
WHERE user_id = $1
  and deleted_at IS  NULL
;

-- name: CreateCart :one
INSERT INTO "shopping_cart"(user_id, goods_id, nums, checked)
VALUES ($1, $2, $3, $4)
returning *;

-- name: DeleteCartItem :one
UPDATE "shopping_cart"
set deleted_at =$1
where user_id = $2
  and goods_id = $3 and deleted_at IS  NULL
returning *;

-- name: UpdateCartItem :one
UPDATE "shopping_cart"
SET updated_at = $1,
    nums       = $2,
    checked    = $3
WHERE user_id = $4
  and goods_id = $5 and deleted_at IS  NULL
returning *;


-- name: GetCartDetailByUIDAndGoodsID :one
SELECT *
FROM "shopping_cart"
where user_id = $1
  and goods_id = $2 and deleted_at IS  NULL;

-- name: GetCartListChecked :many
SELECT *
FROM "shopping_cart"
WHERE user_id = $1 
  and checked = $2 
  and deleted_at IS  NULL;

-- name: CreateOrder :one
INSERT INTO "order_info"(user_id,
                         order_id,
                         status,
                         order_mount,
                         address,
                         signer_name,
                         signer_mobile,
                         post)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetOrderList :many
SELECT *
FROM "order_info"
WHERE user_id = $1 and deleted_at IS  NULL
limit $2 offset $3;


-- name: GetOrderDetail :one
SELECT *
FROM "order_info"
WHERE user_id = $1
  and order_id = $2 and deleted_at IS  NULL
LIMIT 1;


-- name: UpdateOrder :one
update "order_info"
set updated_at = $1,
    pay_type   = $2,
    pay_time   = $3,
    status     = $4
where order_id = $5 and deleted_at IS  NULL
returning *;


-- name: CreateOrderGoods :one
INSERT INTO "order_goods"(ORDER_ID, GOODS_ID, GOODS_NAME, GOODS_PRICE, NUMS)
VALUES ($1, $2, $3, $4, $5)
returning *;

-- name: GetOrderListByOrderID :many
SELECT *
FROM order_goods
WHERE order_id = $1
  and deleted_at IS  NULL;


