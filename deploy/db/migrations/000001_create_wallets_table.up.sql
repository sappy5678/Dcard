BEGIN;
CREATE TABLE short_url (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    expire_time BIGINT NOT NULL,
    created_time BIGINT NOT NULL,
    CONSTRAINT short_url_original_url_check CHECK (original_url <> ''),
    CONSTRAINT short_url_short_code_check CHECK (short_code <> ''),
    CONSTRAINT short_url_expire_time_check CHECK (expire_time > 0),
    CONSTRAINT short_url_created_time_check CHECK (created_time > 0)
);

CREATE INDEX idx_short_code ON short_url (short_code);
COMMIT;
