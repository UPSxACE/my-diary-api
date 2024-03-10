-- name: ListNotes :many
SELECT * FROM note
ORDER BY created_at DESC;

-- name: CreateNote :one
INSERT INTO note(author_id, title, "content", content_raw, created_at)
VALUES($1, $2, $3, $4, NOW())
RETURNING id;