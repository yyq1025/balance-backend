IF OBJECT_ID('wallets', 'U') IS NOT NULL
    DROP TABLE wallets;
IF OBJECT_ID('networks', 'U') IS NOT NULL
    DROP TABLE networks;

CREATE TABLE networks
(
  chain_id VARCHAR(64) NOT NULL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  url VARCHAR(128) NOT NULL,
  symbol VARCHAR(16) NOT NULL,
  explorer VARCHAR(128) NOT NULL
);

CREATE TABLE wallets
(
  id INT IDENTITY(1,1) PRIMARY KEY,
  user_id VARCHAR(128) NOT NULL,
  address VARBINARY(32) NOT NULL,
  network_name VARCHAR(64) NOT NULL,
  token VARBINARY(32) NOT NULL,
  CONSTRAINT UC_Wallets UNIQUE (user_id, address, network_name, token)
);