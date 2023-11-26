-- name: CreateUser :one
INSERT INTO "users" (
  "id",
  "email",
  "hashed_password",
  "line_id",
  "nickname",
  "role"
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "users"
WHERE "email" = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM "users"
LIMIT $1
OFFSET $2;

-- name: UpdateUserByEmail :one
UPDATE "users"
SET
  "hashed_password" = COALESCE(sqlc.narg('hashed_password'), "hashed_password"),
  "line_id" = COALESCE(sqlc.narg('line_id'), "line_id"),
  "nickname" = COALESCE(sqlc.narg('nickname'), "nickname"),
  "is_email_verified" = COALESCE(sqlc.narg('is_email_verified'), "is_email_verified"),
  "updated_at" = COALESCE(sqlc.arg('updated_at'), "updated_at")
WHERE
  "email" = sqlc.arg('email')
RETURNING *;
