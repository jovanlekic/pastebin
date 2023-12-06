-- postgres.down.sql

-- Drop the Object table
DROP TABLE IF EXISTS Object;

-- Drop the Users table
DROP TABLE IF EXISTS Users;

-- Drop the uuid-ossp extension
DROP EXTENSION IF EXISTS "uuid-ossp";
