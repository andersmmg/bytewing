-- Passkeys (WebAuthn), one row per registered authenticator.

-- +goose Up
CREATE TABLE webauthn_credentials (
    id            INTEGER PRIMARY KEY,
    user_id       INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    credential_id BLOB    NOT NULL UNIQUE,
    -- Full serialized credential, including public key and other metadata.
    credential    TEXT    NOT NULL,
    name          TEXT    NOT NULL DEFAULT '',
    created_at    INTEGER NOT NULL,
    last_used_at  INTEGER
) STRICT;

CREATE INDEX webauthn_credentials_user_idx ON webauthn_credentials (user_id);

-- +goose Down
DROP TABLE webauthn_credentials;
