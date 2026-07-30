ALTER TABLE `weapons`
  ADD COLUMN `image_url` VARCHAR(2048) NULL AFTER `element_type`;

ALTER TABLE `items`
  ADD COLUMN `image_url` VARCHAR(2048) NULL AFTER `effect_type`;

ALTER TABLE `monsters`
  ADD COLUMN `image_url` VARCHAR(2048) NULL AFTER `recommended_level`;
