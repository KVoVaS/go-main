-- name: CreateUser :one
INSERT INTO users (name) VALUES ($1) RETURNING id, name;

-- name: AddDocument :one
INSERT INTO documents (title, user_id) VALUES ($1, $2) RETURNING id, title, user_id;

-- name: ListDocumentsByUser :many
SELECT id, title, user_id FROM documents WHERE user_id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);