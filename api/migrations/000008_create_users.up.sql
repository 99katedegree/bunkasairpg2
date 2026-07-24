CREATE TABLE IF NOT EXISTS `users` (
  `id`                  CHAR(36)     NOT NULL,
  `equipped_weapon_id`  BIGINT       NULL,
  `avatar_image_id`     BIGINT       NULL,
  `name`                VARCHAR(255) NOT NULL,
  `level`               INT          NOT NULL DEFAULT 1,
  `hit_point`           INT          NOT NULL,
  `experience_point`    INT          NOT NULL DEFAULT 0,
  `remember_token`      VARCHAR(255) NULL,
  `created_at`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_equipped_weapon_id` FOREIGN KEY (`equipped_weapon_id`) REFERENCES `weapons` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_users_avatar_image_id`    FOREIGN KEY (`avatar_image_id`)    REFERENCES `images`  (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
