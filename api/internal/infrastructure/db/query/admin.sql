-- name: FindAdminByEmail :one
SELECT id, email, password, remember_token, created_at, updated_at
FROM admins
WHERE email = ?
LIMIT 1;
