ALTER TABLE `monsters`
  ADD COLUMN `recommended_level` INT NOT NULL DEFAULT 1 AFTER `experience_point`;
