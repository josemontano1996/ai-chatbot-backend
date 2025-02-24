-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: FindByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;


-- name: FindById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;