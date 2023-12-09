-- postgres.up.sql

-- Ensure the uuid-ossp extension is available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the Users table
CREATE TABLE IF NOT EXISTS Users (
    user_id uuid DEFAULT uuid_generate_v4(),
    name varchar(20) NOT NULL,
    password varchar(32) NOT NULL,
    pasteNum int NOT NULL,
    dev_key varchar(32) NOT NULL,
    email varchar(32) NOT NULL,
    PRIMARY KEY (user_id)
);

-- Create the Object table
CREATE TABLE IF NOT EXISTS Object (
    dev_key varchar(32) NOT NULL,
    paste_key varchar(20) NOT NULL,
    message_id varchar(32),
    PRIMARY KEY (dev_key, paste_key)
);
