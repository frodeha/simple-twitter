CREATE TABLE `Tweets` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `message` text NOT NULL,
  `tag` varchar(32) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `CREATED_AT` (`created_at`) USING BTREE,
  KEY `TAG` (`tag`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;