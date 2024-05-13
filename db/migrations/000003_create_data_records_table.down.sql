BEGIN;

DROP TYPE IF EXISTS DATA_KIND;
DROP TRIGGER IF EXISTS update_data_records_updated_at ON data_records;
DROP TABLE IF EXISTS data_records;

COMMIT;