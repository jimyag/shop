CREATE TABLE "goods"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz          DEFAULT null,
    "name"       varchar     NOT NULL,
    "price"      float       NOT NULL
);

CREATE INDEX ON "goods" ("name");
