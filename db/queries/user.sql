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
SET hashed_password = $2, password_changed_at = now()
WHERE username = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE username = $1
RETURNING *;
