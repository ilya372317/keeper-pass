BEGIN;
DROP TRIGGER IF EXISTS before_update_users_update_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP TABLE users;
COMMIT;