-- +migrate Up

CREATE TABLE wallets(
  id BIGSERIAL NOT NULL  PRIMARY KEY,
  name varchar(150) NOT NULL,
  kind varchar(100) NOT NULL UNIQUE,
);

-- +migrate Down

DROP TABLE wallets;