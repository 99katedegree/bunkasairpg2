CREATE TABLE IF NOT EXISTS `weapons` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT,
  `name`            VARCHAR(255) NOT NULL,
  `index_number`    VARCHAR(255) NOT NULL,
  `physics_attack`  INT          NOT NULL,
  `element_attack`  INT          NULL,
  `physics_type`    VARCHAR(20)  NOT NULL,
  `element_type`    VARCHAR(20)  NOT NULL,
  `created_at`      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_weapons_index_number` (`index_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
