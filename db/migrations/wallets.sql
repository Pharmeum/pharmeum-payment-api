-- +migrate Up

CREATE TABLE wallets(
  public_key varchar(256) NOT NULL  PRIMARY KEY,
  name varchar(150) NOT NULL,
  kind varchar(100) NOT NULL,
  owner_id INTEGER
);

-- +migrate Down

DROP TABLE wallets;