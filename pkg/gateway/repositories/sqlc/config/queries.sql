

-- name: GetAccountByEmail :one
SELECT id, email, password, created_at, updated_at
FROM accounts
WHERE email = $1;

-- name: CreateAccount :one
INSERT INTO accounts (email, password)
VALUES ($1, $2)
RETURNING id;