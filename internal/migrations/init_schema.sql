-- Before applying migrations, do the following in Postgres:
--     1. CREATE ROLE "owner" with login password 'secret';
--     2. CREATE DATABASE shop;
-- +migrate Up
CREATE TABLE accounts(
    ID SERIAL PRIMARY KEY NOT NULL, 
    NAME CHAR(50) NOT NULL, 
    EMAIL CHAR(25) NOT NULL, 
    PASSWORD text NOT NULL,
    ROLE CHAR(5) NOT NULL, 
    IS_ADMIN BOOLEAN NOT NULL
);

CREATE TABLE carts(
    ID SERIAL PRIMARY KEY NOT NULL, 
    ACCOUNT_ID INT NOT NULL, 
    ITEM_ID INT NOT NULL, 
    ITEM_NAME TEXT
);

CREATE TABLE items(
    ID SERIAL PRIMARY KEY NOT NULL, 
    NAME TEXT NOT NULL,
    SKU INT NOT NULL
);

-- +migrate Down
DROP TABLE accounts;
DROP TABLE carts;
DROP TABLE items;