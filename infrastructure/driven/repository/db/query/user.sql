-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: FindByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;


-- name: FindById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users 
SET
email = COALESCE(NULLIF(sqlc.narg(new_email), ''), email),
password = COALESCE(NULLIF(sqlc.narg(new_password), ''), password)
WHERE id = sqlc.arg(id)
RETURNING *;
