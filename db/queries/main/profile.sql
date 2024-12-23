-- name: GetProfile :one
SELECT id, username, email, phone_number, elo, date_joined
FROM profiles
WHERE id = $1;

-- name: InsertProfile :exec
INSERT INTO profiles (username, email, phone_number, elo, date_joined)
VALUES ($1, $2, $3, $4, $5);