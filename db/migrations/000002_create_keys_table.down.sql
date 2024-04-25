BEGIN;
DROP TRIGGER IF EXISTS before_update_keys_update_updated_at ON keys;
DROP TABLE IF EXISTS keys;
COMMIT;