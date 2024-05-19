-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email
)
VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE username = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  password_changed_at = CASE WHEN sqlc.narg(hashed_password) IS NOT NULL
    THEN now() ELSE password_changed_at
    END
WHERE username = sqlc.arg(username)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE username = $1
RETURNING *;
