CREATE TABLE IF NOT EXISTS `monster_entries` (
  `user_id`    CHAR(36)  NOT NULL,
  `monster_id` CHAR(36)  NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `monster_id`),
  CONSTRAINT `fk_monster_entries_user_id`    FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_monster_entries_monster_id` FOREIGN KEY (`monster_id`) REFERENCES `monsters` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
