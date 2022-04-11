-- name: CreateUser :one
INSERT INTO "user"(email,
                   password,
                   nickname,
                   gender,
                   role)
VALUES ($1, $2, $3, $4, $5)
returning *;

-- name: GetUserById :one
SELECT *
FROM "user"
WHERE deleted_at IS NULL
  and id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"
WHERE deleted_at IS NULL
  and email = $1
LIMIT 1;

-- name: ListUsers :many
select *
from "user"
where deleted_at IS NULL
order by id
limit $1 offset $2;


-- name: DeleteUser :execrows
update "user"
set deleted_at =$2
where id = $1
  and deleted_at is null;

-- name: UpdateUser :one
update "user"
set updated_at = $1,
    nickname   = $2,
    gender     = $3,
    role       = $4,
    password   = $5
where id = $6
returning *;

