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

create type GoodsDetail as
(
    goods_id integer,
    nums     integer
);


create table stock_sell_detail
(
    "order_id" int8          not null primary key,
    "status"   int2          not null,
    "detail"   GoodsDetail[] not null
);

CREATE UNIQUE INDEX ON "stock_sell_detail" ("order_id");
