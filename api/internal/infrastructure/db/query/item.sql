-- name: GetAllItemSummaries :many
-- 管理画面の一覧用。図鑑番号は 4 桁ゼロ埋めなので文字列順がそのまま番号順。
SELECT id, name, index_number FROM items ORDER BY index_number;

-- name: GetItemByID :one
SELECT i.id, i.name, i.index_number, i.effect_type, i.created_at, i.updated_at,
    hi.amount, bi.rate as buff_rate, bi.target as buff_target, di.rate as debuff_rate, di.target as debuff_target
FROM items i
LEFT JOIN heal_items hi ON i.id = hi.item_id
LEFT JOIN buff_items bi ON i.id = bi.item_id
LEFT JOIN debuff_items di ON i.id = di.item_id
WHERE i.id = ? LIMIT 1;

-- name: GetItemsByUserID :many
SELECT i.id, i.name, i.index_number, i.effect_type, i.created_at, i.updated_at,
    hi.amount, bi.rate as buff_rate, bi.target as buff_target, di.rate as debuff_rate, di.target as debuff_target, ui.count
FROM items i
INNER JOIN user_items ui ON i.id = ui.item_id
LEFT JOIN heal_items hi ON i.id = hi.item_id
LEFT JOIN buff_items bi ON i.id = bi.item_id
LEFT JOIN debuff_items di ON i.id = di.item_id
WHERE ui.user_id = ?;

-- name: GetItemIndexByUserID :many
SELECT i.id, i.name, i.index_number, i.effect_type, i.created_at, i.updated_at,
    hi.amount, bi.rate as buff_rate, bi.target as buff_target, di.rate as debuff_rate, di.target as debuff_target
FROM items i
INNER JOIN item_entries ie ON i.id = ie.item_id
LEFT JOIN heal_items hi ON i.id = hi.item_id
LEFT JOIN buff_items bi ON i.id = bi.item_id
LEFT JOIN debuff_items di ON i.id = di.item_id
WHERE ie.user_id = ?
ORDER BY i.id
LIMIT ? OFFSET ?;

-- name: CountItemIndexByUserID :one
SELECT COUNT(*) FROM item_entries WHERE user_id = ?;

-- name: DecrementUserItem :exec
UPDATE user_items SET count = count - 1, updated_at = NOW()
WHERE user_id = ? AND item_id = ? AND count > 0;

-- name: DeleteUserItemIfZero :exec
DELETE FROM user_items WHERE user_id = ? AND item_id = ? AND count = 0;

-- name: GrantItemToUser :exec
INSERT INTO user_items (user_id, item_id, count, created_at, updated_at)
VALUES (?, ?, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE count = count + 1, updated_at = NOW();

-- name: GrantItemToEntry :exec
INSERT IGNORE INTO item_entries (user_id, item_id, created_at, updated_at)
VALUES (?, ?, NOW(), NOW());
