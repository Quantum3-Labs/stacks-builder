-- Disable foreign key constraints before dropping
PRAGMA foreign_keys = OFF;

DROP INDEX IF EXISTS idx_query_logs_endpoint;
DROP INDEX IF EXISTS idx_query_logs_created_at;
DROP INDEX IF EXISTS idx_query_logs_user_id;
DROP TABLE IF EXISTS query_logs;

-- Re-enable foreign key constraints
PRAGMA foreign_keys = ON;
