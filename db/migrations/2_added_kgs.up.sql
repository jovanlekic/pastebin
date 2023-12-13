-- postgres.up.sql

-- Create the Used table
CREATE TABLE IF NOT EXISTS Used (
    id uuid DEFAULT uuid_generate_v4(),
    key varchar(32) NOT NULL,
    PRIMARY KEY (id)
);

-- Create the Unused table
CREATE TABLE IF NOT EXISTS Unused (
    id uuid DEFAULT uuid_generate_v4(),
    key varchar(32) NOT NULL,
    PRIMARY KEY (id)
);

-- Create the Migration table
CREATE TABLE IF NOT EXISTS Migration (
    id uuid DEFAULT uuid_generate_v4(),
    migration varchar(32) NOT NULL,
    PRIMARY KEY (id)
);
