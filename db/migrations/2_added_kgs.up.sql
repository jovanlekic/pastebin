-- postgres.up.sql

-- Create the Keys table
CREATE TABLE IF NOT EXISTS Keys (
    id uuid DEFAULT uuid_generate_v4(),
    key varchar(32) NOT NULL,
    used boolean DEFAULT false, 
    PRIMARY KEY (id)
);


-- Create the Migration table
CREATE TABLE IF NOT EXISTS Migration (
    id uuid DEFAULT uuid_generate_v4(),
    migration varchar(32) NOT NULL,
    PRIMARY KEY (id)
);
