-- name: GetWebAuthnCredentialByCredentialID :one
-- The login lookup: an assertion carries the credential id and nothing else
-- until it has been verified.
SELECT * FROM webauthn_credentials WHERE credential_id = ?;

-- name: ListWebAuthnCredentialsForUser :many
-- The management list in Settings. The stored credential is never sent to the
-- client.
SELECT * FROM webauthn_credentials
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CountWebAuthnCredentials :one
-- Whether the instance has any passkey at all, so the sign-in screen can offer
-- one only when signing in with it can work.
SELECT COUNT(*) FROM webauthn_credentials;

-- name: CreateWebAuthnCredential :one
INSERT INTO webauthn_credentials (user_id, credential_id, credential, name, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateWebAuthnCredentialUse :exec
-- Records a successful assertion: the credential as the verification left it
-- (its signature counter is what detects a cloned authenticator) and the use
-- time for the Settings list. Best-effort, like the session touch: a lost
-- timestamp must never fail an otherwise-valid sign-in.
UPDATE webauthn_credentials
SET credential = ?, last_used_at = ?
WHERE id = ?;

-- name: DeleteWebAuthnCredentialForUser :execrows
-- Scoped to the owner so one account cannot revoke another's passkey by
-- guessing an id. execrows lets the handler answer 404 when nothing matched.
DELETE FROM webauthn_credentials
WHERE id = ? AND user_id = ?;
