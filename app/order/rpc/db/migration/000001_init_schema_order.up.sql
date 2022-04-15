CREATE TABLE "shopping_cart"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz          DEFAULT null,
    "user_id"    integer     NOT NULL,
    "goods_id"   integer     NOT NULL,
    "nums"       integer     NOT NULL,
    "checked"    boolean     NOT NULL
);

CREATE TABLE "order_info"
(
    "id"            bigserial PRIMARY KEY,
    "created_at"    timestamptz NOT NULL DEFAULT (now()),
    "updated_at"    timestamptz NOT NULL DEFAULT (now()),
    "deleted_at"    timestamptz          DEFAULT null,
    "user_id"       integer     NOT NULL,
    "order_id"      int8 UNIQUE NOT NULL,
    "pay_type"      varchar,
    "status"        int2        NOT NULL, -- 1 待支付 2 成功 3 超时关闭
    "trade_id"      varchar,              --支付编号
    "order_mount"   float,                -- 订单金额
    "pay_time"      timestamptz,
    "address"       varchar     NOT NULL,
    "signer_name"   varchar(40) NOT NULL,
    "signer_mobile" varchar(20) NOT NULL,
    "post"          varchar     NOT NULL
);

CREATE TABLE "order_goods"
(
    "id"          bigserial PRIMARY KEY,
    "created_at"  timestamptz NOT NULL DEFAULT (now()),
    "updated_at"  timestamptz NOT NULL DEFAULT (now()),
    "deleted_at"  timestamptz          DEFAULT null,
    "order_id"    int8        NOT NULL,
    "goods_id"    integer     NOT NULL,
    "goods_name"  varchar     NOT NULL,
    "goods_price" float       NOT NULL,
    "nums"        integer     NOT NULL
);

CREATE INDEX ON "shopping_cart" ("user_id");

CREATE INDEX ON "shopping_cart" ("goods_id");

CREATE INDEX ON "order_info" ("user_id");

CREATE INDEX ON "order_info" ("order_id");

CREATE INDEX ON "order_goods" ("order_id");

CREATE INDEX ON "order_goods" ("goods_id");

CREATE INDEX ON "order_goods" ("goods_name");
