-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password,
    host
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;
