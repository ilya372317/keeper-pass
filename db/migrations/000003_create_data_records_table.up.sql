BEGIN;

CREATE TABLE data_records
(
    id         SERIAL PRIMARY KEY,
    payload    JSON         NOT NULL,
    metadata   JSON         NOT NULL,
    nonce      VARCHAR(255) NOT NULL,
    crypto_key VARCHAR(255) NOT NULL,
    kind       SMALLINT     NOT NULL,
    user_id    INTEGER      NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_data_records_updated_at
    BEFORE UPDATE
    ON data_records
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMENT ON COLUMN data_records.id IS 'Идентификатор';
COMMENT ON COLUMN data_records.payload IS 'Данные';
COMMENT ON COLUMN data_records.metadata IS 'Метаданные';
COMMENT ON COLUMN data_records.nonce IS 'Уникальная последовательность для шифрования и дешифрования';
COMMENT ON COLUMN data_records.crypto_key IS 'Ключ для шифрования';
COMMENT ON COLUMN data_records.kind IS 'Тип данных';
COMMENT ON COLUMN data_records.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN data_records.updated_at IS 'Дата обновления';
COMMENT ON COLUMN data_records.created_at IS 'Дата создания';

COMMIT;