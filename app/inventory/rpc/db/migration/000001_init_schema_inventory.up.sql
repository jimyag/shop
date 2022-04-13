CREATE TABLE "inventory"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz          DEFAULT null,
    "goods_id"   integer     NOT NULL,
    "sticks"     integer     NOT NULL,
    "version"    integer     NOT NULL
);

CREATE INDEX ON "inventory" ("goods_id");

CREATE tABLE "stock_sell_detail"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz          DEFAULT null,
    "order_id"   integer     NOT NULL, -- index
    "status"     integer     NOT NULL, -- 1 表示已经扣减 2 已归还
    "detail"     integer     NOT NULL  -- 包含了 good_id 和 num
)
