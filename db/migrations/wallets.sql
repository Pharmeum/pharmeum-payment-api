-- +migrate Up

CREATE TABLE wallets(
  public_key varchar(256) NOT NULL  PRIMARY KEY,
  owner_id INTEGER
);

-- +migrate Down

DROP TABLE wallets;