DROP TABLE IF EXISTS `wallets`;
DROP TABLE IF EXISTS `networks`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` INT AUTO_INCREMENT NOT NULL,
  `email` VARCHAR(320) NOT NULL,
  `password` BINARY(60) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`email`)
);

CREATE TABLE `networks` (
  `name` VARCHAR(64) NOT NULL,
  `url` VARCHAR(128) NOT NULL,
  `symbol` VARCHAR(16) NOT NULL,
  PRIMARY KEY (`name`)
);

CREATE TABLE `wallets` (
  `id` INT AUTO_INCREMENT NOT NULL,
  `user_id` INT NOT NULL,
  `address` BINARY(20) NOT NULL,
  `network` VARCHAR(64) NOT NULL,
  `token` BINARY(20) NOT NULL,
  `tag` VARCHAR(320) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`user_id`, `address`, `network`, `token`, `tag`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`network`) REFERENCES networks(`name`) ON DELETE CASCADE
);