DROP TABLE IF EXISTS "wallets";
DROP TABLE IF EXISTS "networks";
DROP TABLE IF EXISTS "users";

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "email" VARCHAR(320) NOT NULL,
  "password" BYTEA NOT NULL,
  UNIQUE ("email")
);

CREATE TABLE "networks" (
  "name" VARCHAR(64) NOT NULL PRIMARY KEY,
  "url" VARCHAR(128) NOT NULL,
  "symbol" VARCHAR(16) NOT NULL
);

CREATE TABLE "wallets" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "address" BYTEA NOT NULL,
  "network" VARCHAR(64) NOT NULL,
  "token" BYTEA NOT NULL,
  "tag" VARCHAR(255) NOT NULL,
  UNIQUE ("user_id", "address", "network", "token", "tag"),
  FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE,
  FOREIGN KEY ("network") REFERENCES networks("name") ON DELETE CASCADE
);