-- name: ListUser :many
SELECT * FROM "user" ORDER BY id ASC;
-- name: CreateUser :one
INSERT INTO public."user"(username, password, email, avatar_url, full_name, created_at, role_id)
VALUES($1, $2, $3, $4, $5, NOW(), $6)
RETURNING id;
-- name: GetUserAuthByUsername :one
SELECT "user".id, password, role_id, can_all FROM "user" 
INNER JOIN role ON "user".role_id = role.id
WHERE username = $1;