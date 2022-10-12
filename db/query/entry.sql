-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEmtry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEmtries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;