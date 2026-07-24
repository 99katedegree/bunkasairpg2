-- name: CreateImage :exec
INSERT INTO images (directory, url, created_at)
VALUES (?, ?, NOW());

-- name: GetLastInsertImage :one
SELECT * FROM images WHERE id = LAST_INSERT_ID() LIMIT 1;
