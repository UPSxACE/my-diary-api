-- name: ListNotes :many
SELECT * FROM note
ORDER BY created_at DESC;