-- Disable foreign key constraints before dropping tables
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS ingestion_jobs;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS users;

-- Re-enable foreign key constraints
PRAGMA foreign_keys = ON;
