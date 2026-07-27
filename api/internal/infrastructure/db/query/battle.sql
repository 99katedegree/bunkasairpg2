-- name: CreateBattle :exec
INSERT INTO battles (id, user_id, monster_id, start_weapon_id, seed, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetBattleByID :one
SELECT * FROM battles WHERE id = ? LIMIT 1;

-- name: UpdateBattleStatus :exec
UPDATE battles SET status = ?, updated_at = NOW() WHERE id = ?;

-- name: UpsertBossRecord :exec
INSERT INTO boss_records (user_id, clear_time, created_at, updated_at)
VALUES (?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    clear_time = IF(VALUES(clear_time) < clear_time, VALUES(clear_time), clear_time),
    updated_at = NOW();
