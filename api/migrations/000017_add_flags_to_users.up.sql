ALTER TABLE `users`
  ADD COLUMN `is_archived`  TINYINT(1) NOT NULL DEFAULT 0 AFTER `experience_point`,
  ADD COLUMN `is_activated` TINYINT(1) NOT NULL DEFAULT 0 AFTER `is_archived`;
