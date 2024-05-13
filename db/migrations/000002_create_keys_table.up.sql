BEGIN;

CREATE TABLE keys
(
    id         serial PRIMARY KEY,
    key        varchar(255)              NOT NULL,
    is_current BOOLEAN                   NOT NULL DEFAULT true,
    nonce varchar(255) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL
);

COMMENT ON column keys.id IS 'identifier';
COMMENT ON column keys.key IS 'key for crypt and decrypt data keys';
COMMENT ON column keys.is_current IS 'flag indicates is current key or not';
COMMENT ON column keys.created_at IS 'creation date';
COMMENT ON column keys.updated_at IS 'last update date';

CREATE TRIGGER before_update_keys_update_updated_at
    BEFORE UPDATE
    ON keys
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;