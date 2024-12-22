-- name: CreateReceivable :one
INSERT INTO receivables (id, created_at, updated_at, userid, date, amount_total, amount_received, amount_left, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteReceivable :one
DELETE FROM receivables WHERE id=$1
RETURNING *;

-- name: GetAllReceivables :many
SELECT r.*, u.name AS username FROM receivables r
    JOIN users u ON r.userid = u.id;

-- name: UpdateReceivable :one
UPDATE receivables
SET
    date=$2,
    updated_at = NOW(),
    amount_total = $3,
    amount_received = $4,
    amount_left = $5
WHERE id=$1
RETURNING *;

-- name: GetReceivableByID :one
SELECT r.*, u.name AS username FROM receivables r
           JOIN users u ON r.userid = u.id
           WHERE r.id = $1;

-- name: GetReceivablesByUserId :many
SELECT r.*, u.name AS username FROM receivables r
           JOIN users u ON r.userid = u.id
           WHERE r.userid = $1;
