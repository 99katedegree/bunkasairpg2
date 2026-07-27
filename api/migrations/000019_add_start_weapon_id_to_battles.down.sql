ALTER TABLE `battles`
  DROP FOREIGN KEY `fk_battles_start_weapon_id`,
  DROP COLUMN `start_weapon_id`;
