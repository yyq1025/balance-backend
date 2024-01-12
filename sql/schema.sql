DROP TABLE IF EXISTS `wallets`;
DROP TABLE IF EXISTS `networks`;

CREATE TABLE `networks` (
  `chain_id` VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  url VARCHAR(128) NOT NULL,
  symbol VARCHAR(16) NOT NULL,
  explorer VARCHAR(128) NOT NULL,
  PRIMARY KEY (name)
);

CREATE TABLE wallets (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(128) NOT NULL,
  address BLOB NOT NULL,
  network_name VARCHAR(64) NOT NULL,
  token BLOB NOT NULL,
  UNIQUE (user_id, address(255), network_name, token(255))
);