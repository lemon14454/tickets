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
select * from users
where email = $1 limit 1;
