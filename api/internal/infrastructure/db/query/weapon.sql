-- name: GetAllWeaponIDs :many
SELECT id FROM weapons ORDER BY id;

-- name: GetWeaponByID :one
SELECT * FROM weapons WHERE id = ? LIMIT 1;

-- name: GetWeaponsByUserID :many
SELECT w.* FROM weapons w
INNER JOIN user_weapons uw ON w.id = uw.weapon_id
WHERE uw.user_id = ?;

-- name: GetWeaponIndexByUserID :many
SELECT w.* FROM weapons w
INNER JOIN weapon_entries we ON w.id = we.weapon_id
WHERE we.user_id = ?
ORDER BY w.id
LIMIT ? OFFSET ?;

-- name: CountWeaponIndexByUserID :one
SELECT COUNT(*) FROM weapon_entries WHERE user_id = ?;

-- name: IsWeaponOwnedByUser :one
SELECT COUNT(*) FROM user_weapons WHERE user_id = ? AND weapon_id = ?;

-- name: GrantWeaponToUser :exec
INSERT IGNORE INTO user_weapons (user_id, weapon_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW());

-- name: GrantWeaponToEntry :exec
INSERT IGNORE INTO weapon_entries (user_id, weapon_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW());
