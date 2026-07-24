-- name: GetUserByID :one
SELECT u.id, u.equipped_weapon_id, u.avatar_image_id, u.name, u.level, u.hit_point, u.experience_point, u.is_archived, u.is_activated, u.remember_token, u.created_at, u.updated_at,
    w.id as weapon_id_join, w.name as weapon_name, w.index_number as weapon_index_number, w.physics_attack, w.element_attack, w.physics_type, w.element_type,
    w.created_at as weapon_created_at, w.updated_at as weapon_updated_at,
    i.url as avatar_image_url
FROM users u
LEFT JOIN weapons w ON u.equipped_weapon_id = w.id
LEFT JOIN images i ON u.avatar_image_id = i.id
WHERE u.id = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users SET
    name = COALESCE(?, name),
    level = COALESCE(?, level),
    hit_point = COALESCE(?, hit_point),
    experience_point = COALESCE(?, experience_point),
    equipped_weapon_id = ?,
    avatar_image_id = ?,
    updated_at = NOW()
WHERE id = ?;

-- name: ActivateUser :exec
UPDATE users SET is_activated = 1, updated_at = NOW() WHERE id = ?;

-- name: ArchiveAllUsers :exec
UPDATE users SET is_archived = 1, updated_at = NOW() WHERE is_archived = 0;

-- name: DeleteInactiveUsers :exec
DELETE FROM users WHERE is_activated = 0;

-- name: CreateUser :exec
INSERT INTO users (id, name, level, hit_point, experience_point, is_archived, is_activated, created_at, updated_at)
VALUES (?, '名無し', 1, 100, 0, 0, 0, NOW(), NOW());

-- name: GetAllUserIDs :many
SELECT id FROM users WHERE is_archived = 0 ORDER BY created_at;
