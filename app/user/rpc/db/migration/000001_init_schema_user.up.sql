CREATE TABLE "user"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz    NOT NULL DEFAULT (now()),
    "updated_at" timestamptz    NOT NULL DEFAULT (now()),
    "deleted_at" timestamptz             DEFAULT null,
    "email"      varchar UNIQUE NOT NULL,
    "password"   varchar        NOT NULL,
    "nickname"   varchar        NOT NULL,
    "gender"     varchar(6)     NOT NULL DEFAULT 'male',
    "role"       int8           NOT NULL DEFAULT 1
);

CREATE INDEX ON "user" ("email");

COMMENT ON COLUMN "user"."email" IS 'user email';

COMMENT ON COLUMN "user"."password" IS 'user password';

COMMENT ON COLUMN "user"."nickname" IS 'user nickname default email';

COMMENT ON COLUMN "user"."gender" IS 'male man ,female women';

COMMENT ON COLUMN "user"."role" IS '1 user 2 admin';
