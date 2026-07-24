-- name: GetMonsterByID :one
SELECT m.id, m.weapon_id, m.item_id, m.index_number, m.name, m.attack, m.hit_point, m.experience_point,
    m.slash, m.blow, m.shoot, m.neutral, m.flame, m.water, m.wood, m.shine, m.dark,
    m.created_at, m.updated_at,
    w.id as weapon_id_j, w.name as weapon_name, w.index_number as weapon_index_number,
    w.physics_attack, w.element_attack, w.physics_type, w.element_type,
    w.created_at as weapon_created_at, w.updated_at as weapon_updated_at,
    i.id as item_id_j, i.name as item_name, i.index_number as item_index_number, i.effect_type,
    i.created_at as item_created_at, i.updated_at as item_updated_at,
    hi.amount, bi.rate as buff_rate, bi.target as buff_target, di.rate as debuff_rate, di.target as debuff_target
FROM monsters m
LEFT JOIN weapons w ON m.weapon_id = w.id
LEFT JOIN items i ON m.item_id = i.id
LEFT JOIN heal_items hi ON i.id = hi.item_id
LEFT JOIN buff_items bi ON i.id = bi.item_id
LEFT JOIN debuff_items di ON i.id = di.item_id
WHERE m.id = ? LIMIT 1;

-- name: GetMonsterCatalogByUserID :many
SELECT me.monster_id FROM monster_entries me
WHERE me.user_id = ?
ORDER BY me.created_at
LIMIT ? OFFSET ?;

-- name: CountMonsterCatalogByUserID :one
SELECT COUNT(*) FROM monster_entries WHERE user_id = ?;

-- name: RegisterMonsterEntry :exec
INSERT IGNORE INTO monster_entries (user_id, monster_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW());

-- name: IsMonsterEntryRegistered :one
SELECT COUNT(*) FROM monster_entries WHERE user_id = ? AND monster_id = ?;
