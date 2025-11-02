-- Initialization script for debug database
-- This script runs only on first container start

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create test data (опционально)
-- INSERT INTO accounts (id, username, email) VALUES
--   ('test-user-1', 'debuguser', 'debug@example.com');

-- Grant all privileges
GRANT ALL PRIVILEGES ON DATABASE wishlist TO wishlist_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO wishlist_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO wishlist_user;

-- Create useful views for debugging
CREATE OR REPLACE VIEW v_table_sizes AS
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Log message
DO $$
BEGIN
    RAISE NOTICE 'Debug database initialized successfully!';
END $$;

