DROP TABLE IF EXISTS "wallets";
DROP TABLE IF EXISTS "networks";

CREATE TABLE "networks" (
  "chain_id" VARCHAR(64) NOT NULL,
  "name" VARCHAR(64) NOT NULL PRIMARY KEY,
  "url" VARCHAR(128) NOT NULL,
  "symbol" VARCHAR(16) NOT NULL,
  "explorer" VARCHAR(128) NOT NULL
);

CREATE TABLE "wallets" (
  "id" SERIAL PRIMARY KEY,
  "user_id" VARCHAR(128) NOT NULL,
  "address" BYTEA NOT NULL,
  "network_name" VARCHAR(64) NOT NULL,
  "token" BYTEA NOT NULL,
  UNIQUE ("user_id", "address", "network_name", "token"),
  FOREIGN KEY ("network_name") REFERENCES networks("name") ON DELETE CASCADE
);