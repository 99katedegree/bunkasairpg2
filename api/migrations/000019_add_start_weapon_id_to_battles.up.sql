-- バトル開始時に装備していた武器を記録する。
-- 戦闘中の武器変更は users.equipped_weapon_id を書き換えてしまうため、
-- 終了時に再計算しようとすると「最後に持ち替えた武器」から始めることになり、
-- クライアントの結果と一致しなくなる。開始時点を固定するための列。
ALTER TABLE `battles`
  ADD COLUMN `start_weapon_id` BIGINT NULL AFTER `monster_id`,
  ADD CONSTRAINT `fk_battles_start_weapon_id`
    FOREIGN KEY (`start_weapon_id`) REFERENCES `weapons` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;
