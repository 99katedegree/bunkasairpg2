CREATE TABLE IF NOT EXISTS `battles` (
  `id`         CHAR(36)    NOT NULL,
  `user_id`    CHAR(36)    NOT NULL,
  `monster_id` CHAR(36)    NULL,
  `seed`       BIGINT      NOT NULL,
  `status`     VARCHAR(20) NOT NULL COMMENT 'in_progress|completed|lost|expired',
  `created_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_battles_user_id`    FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_battles_monster_id` FOREIGN KEY (`monster_id`) REFERENCES `monsters` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
