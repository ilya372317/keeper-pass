BEGIN;

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id              SERIAL PRIMARY KEY,
    hashed_password text                      NOT NULL,
    email           text                      NOT NULL,
    salt varchar(255) NOT NULL,
    created_at      timestamptz DEFAULT NOW() NOT NULL,
    updated_at      timestamptz DEFAULT NOW() NOT NULL
);

COMMENT ON column users.id IS 'identifier';
COMMENT ON column users.hashed_password IS 'user hashed password';
COMMENT ON column users.email IS 'user email';
COMMENT ON column users.salt IS 'random generated salt for hash check';


CREATE TRIGGER before_update_users_update_updated_at
    BEFORE update
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
