# API

Everything the interface does is a plain JSON API under `/api/v1`. This page
lists the endpoints and how to authenticate. There is no generated OpenAPI
contract yet, but the API is stable and usable directly.

Everything is authenticated unless it appears under
[Setup, invitations and recovery](#setup-invitations-and-recovery), which lists
the public endpoints.

## Authentication

Three ways in resolve to the same account:

- **Session cookie**: a browser signs in via `POST /api/v1/auth/login` and
  receives an `HttpOnly` cookie. This is what the web UI uses.
- **Passkey**: a browser asks `POST /api/v1/auth/passkeys/login/options` for a
  challenge, answers it with the authenticator, and posts the assertion to
  `POST /api/v1/auth/passkeys/login/verify`, which sets the same cookie.
- **Bearer token**: a program sends an [API token](/reference/api-tokens) as
  `Authorization: Bearer bytewing_live_…`. This is what scripts and integrations
  use.

A token acts with the full authority of the account that created it. Every path
traverses exactly the same permission checks.

```sh
TOKEN=bytewing_live_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
curl -H "Authorization: Bearer $TOKEN" \
  "https://bytewing.example.tld/api/v1/fs?path=/"
```

## Errors

Failures return an [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) problem
document with a machine-readable `code`:

```json
{
  "type": "…",
  "title": "Not signed in",
  "status": 401,
  "code": "unauthenticated",
  "detail": "…"
}
```

Field-level validation errors add an `errors` map of field name to message.

## Setup, invitations and recovery

These are the only endpoints that work without a session or token. They are how
an account gets created or regained in an instance that never sends email.

| Method | Path                            | Description                                                       |
| ------ | ------------------------------- | ----------------------------------------------------------------- |
| `GET`  | `/healthz`                      | Liveness check                                                    |
| `GET`  | `/api/v1/setup`                 | `{ "needs_setup": true }` while no account exists                 |
| `POST` | `/api/v1/setup`                 | Create the first administrator: `{ "token", "username", "password" }`; returns a session token, sets the cookie, and returns the recovery codes |
| `GET`  | `/api/v1/auth/passkeys/status`  | Whether the instance has any passkeys at all                      |
| `POST` | `/api/v1/auth/reset`            | Reset a forgotten password: `{ "username", "code", "new_password" }`; ends every session |
| `GET`  | `/api/v1/invitations/check`     | Validate an invitation token (`?token=`)                          |
| `POST` | `/api/v1/invitations/accept`    | Accept it: `{ "token", "username", "password" }`                 |

The setup token is printed in the server log on first start and expires after 24
hours.

## Authentication & account

| Method   | Path                                       | Description                                                                 |
| -------- | ------------------------------------------ | --------------------------------------------------------------------------- |
| `POST`   | `/api/v1/auth/login`                       | Sign in: `{ "username", "password", "totp_code?" }`; sets the session cookie |
| `POST`   | `/api/v1/auth/passkeys/login/options`      | Start a passkey sign-in; body `{ "username?": "…" }`                        |
| `POST`   | `/api/v1/auth/passkeys/login/verify`       | Finish it: `{ "ceremony_id", "credential" }`; sets the cookie              |
| `POST`   | `/api/v1/auth/logout`                      | Revoke the current session                                                  |
| `GET`    | `/api/v1/auth/me`                          | The signed-in account                                                       |
| `PUT`    | `/api/v1/auth/theme`                       | Set the interface colour theme `{ "theme": "pine" }`; `""` resets           |
| `POST`   | `/api/v1/auth/password`                    | Change password: `{ "current_password", "new_password" }`; ends other sessions |
| `GET`    | `/api/v1/auth/sessions`                    | List active sessions, with device and location                              |
| `POST`   | `/api/v1/auth/sessions/revoke-others`      | Sign out every other device                                                 |
| `DELETE` | `/api/v1/auth/sessions/{id}`               | End one session                                                             |
| `POST`   | `/api/v1/auth/totp/enroll`                 | Mint a candidate TOTP secret and its `otpauth://` URI                       |
| `POST`   | `/api/v1/auth/totp/enable`                 | Confirm it: `{ "secret", "code" }`                                         |
| `POST`   | `/api/v1/auth/totp/disable`                | Turn two-factor off: `{ "code" }`                                          |
| `GET`    | `/api/v1/auth/passkeys`                    | List the passkeys on the account                                            |
| `POST`   | `/api/v1/auth/passkeys/register/options`   | Start registering: `{ "password" }`                                       |
| `POST`   | `/api/v1/auth/passkeys/register/verify`    | Finish it: `{ "ceremony_id", "name", "credential" }`                       |
| `DELETE` | `/api/v1/auth/passkeys/{id}`               | Remove a passkey: `{ "password" }`                                        |
| `GET`    | `/api/v1/auth/recovery`                    | How many unused recovery codes remain                                       |
| `POST`   | `/api/v1/auth/recovery`                    | Replace the set; the new codes are returned once                            |

## Files

| Method   | Path                        | Description                                                                                                          |
| -------- | --------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `GET`    | `/api/v1/fs?path=`          | List a directory                                                                                                     |
| `GET`    | `/api/v1/fs/stat?path=`     | Metadata for one entry                                                                                               |
| `GET`    | `/api/v1/fs/info?path=`     | Size, dates and owner for one entry                                                                                  |
| `GET`    | `/api/v1/fs/count?path=`    | Files and folders count directly inside a folder                                                                     |
| `GET`    | `/api/v1/fs/search?q=`      | Recursive search (only what you may see): `limit`, `path`                                                           |
| `POST`   | `/api/v1/fs/dirs`           | Create a directory: `{ "path": "/x" }`                                                                              |
| `POST`   | `/api/v1/fs/move`           | Move or rename                                                                                                       |
| `POST`   | `/api/v1/fs/copy`           | Copy (a tree or an oversized file runs as a background job)                                                          |
| `POST`   | `/api/v1/fs/extract`        | Extract a ZIP into a new folder (background job): `{ "archive": "/x.zip" }`; not available on object storage          |
| `POST`   | `/api/v1/fs/fetch`          | Download a URL into a folder (background job): `{ "url": "https://…", "dir": "/" }`; needs `BYTEWING_ALLOW_URL_FETCH` |
| `DELETE` | `/api/v1/fs?path=`          | Move to trash                                                                                                        |
| `GET`    | `/api/v1/fs/link?path=`     | Mint a short-lived signed download URL                                                                               |
| `GET`    | `/api/v1/fs/text?path=`     | Read a text file (capped at 2 MiB by default)                                                                        |
| `POST`   | `/api/v1/fs/bundle`         | Mint a signed zip link for a selection or folder                                                                     |
| `GET`    | `/api/v1/fs/thumb?path=`    | Image/video thumbnail: `s` for its pixel size                                                                       |
| `GET`    | `/api/v1/fs/space`          | Free space and read-only state                                                                                       |
| `GET`    | `/api/v1/fs/checksum?path=` | SHA-256 of a file (cached, or a job to poll)                                                                         |
| `PUT`    | `/api/v1/fs/color`          | Set a per-user label colour: `{ "path", "color" }`; `""` clears it                                                  |
| `PUT`    | `/api/v1/fs/reference`      | Make an entry a per-user reference to another: `{ "path", "target" }`                                               |
| `DELETE` | `/api/v1/fs/reference`      | Drop the reference at `?path=`                                                                                       |
| `GET`    | `/api/v1/config`            | Instance capabilities and version                                                                                    |

Labels and references are per-user conveniences: they change nothing about who
may read a path. A reference must point at a file the caller can read when it is
created.

## Sharing

| Method   | Path                  | Description                                        |
| -------- | --------------------- | -------------------------------------------------- |
| `GET`    | `/api/v1/shares`      | List your share links                              |
| `POST`   | `/api/v1/shares`      | Create one: `{ "path", "expires_in_hours?", "password?" }`; a password only applies to a file |
| `DELETE` | `/api/v1/shares/{id}` | Revoke a share link                                |

## Trash

| Method   | Path                         | Description                  |
| -------- | ---------------------------- | ---------------------------- |
| `GET`    | `/api/v1/trash`              | List trashed entries         |
| `POST`   | `/api/v1/trash/{id}/restore` | Restore an entry             |
| `DELETE` | `/api/v1/trash/{id}`         | Delete one entry permanently |
| `DELETE` | `/api/v1/trash`              | Empty the trash              |

## Background jobs

| Method | Path                       | Description                                  |
| ------ | -------------------------- | -------------------------------------------- |
| `GET`  | `/api/v1/jobs`             | List background jobs (copies, etc.)          |
| `GET`  | `/api/v1/jobs/{id}`        | One job's status, progress and byte counts   |
| `POST` | `/api/v1/jobs/{id}/cancel` | Cancel a job (discards a download's partial) |
| `POST` | `/api/v1/jobs/{id}/pause`  | Pause a running download                     |
| `POST` | `/api/v1/jobs/{id}/resume` | Resume a paused download                     |

A job carries `bytes_done` and `bytes_total`; a client shows a transfer rate by
dividing the change in `bytes_done` by the time between polls. Pause and resume
apply to downloads: a resumed download continues from where it stopped using an
HTTP range request, or restarts cleanly if the source does not support ranges.

## API tokens

| Method   | Path                  | Description                    |
| -------- | --------------------- | ------------------------------ |
| `GET`    | `/api/v1/tokens`      | List your tokens (prefix only) |
| `POST`   | `/api/v1/tokens`      | Create a token (returned once) |
| `DELETE` | `/api/v1/tokens/{id}` | Revoke a token                 |

See [API tokens](/reference/api-tokens) for the full flow.

## Administration

These require an administrator account.

| Method         | Path                                            | Description                           |
| -------------- | ----------------------------------------------- | ------------------------------------- |
| `GET`          | `/api/v1/users`                                 | List accounts                         |
| `PATCH`        | `/api/v1/users/{id}`                            | Promote, disable or update an account |
| `DELETE`       | `/api/v1/users/{id}`                            | Remove an account                     |
| `GET`          | `/api/v1/invitations` · `POST` · `DELETE /{id}` | Manage invitations                    |
| `GET`          | `/api/v1/audit`                                 | Read the audit log (keyset paginated) |
| `GET`          | `/api/v1/groups` · `POST` · `DELETE /{id}`      | Manage groups                         |
| `GET`          | `/api/v1/groups/{id}/members`                   | List a group's members                |
| `PUT`/`DELETE` | `/api/v1/groups/{id}/members/{userID}`          | Group membership                      |
| `GET`          | `/api/v1/settings/retention`                    | Retention policy in days              |
| `PUT`          | `/api/v1/settings/retention`                    | Set it: `{ "audit_days", "trash_days" }` |

## Permissions

| Method   | Path                        | Description                     |
| -------- | --------------------------- | ------------------------------- |
| `GET`    | `/api/v1/permissions?path=` | Effective permissions at a path |
| `GET`    | `/api/v1/access?path=`      | Access rules at a path          |
| `POST`   | `/api/v1/access`            | Grant access to a user or group |
| `DELETE` | `/api/v1/access/{id}`       | Revoke an access rule           |
| `GET`    | `/api/v1/shared-with-me`    | Folders granted to you, for a member's sidebar |

## Uploads

Uploads are resumable. The protocol follows [tus](https://tus.io) for the streamed
mode: `Upload-Length` carries the size, `Upload-Metadata` carries base64
`path` and an optional `checksum`, and `PATCH` carries `Upload-Offset`.

The create response includes a `mode`. In `tus` mode bytes go to the server
through `PATCH`; in `direct` mode (object storage) the browser uploads parts
straight to the bucket, so `PATCH` is refused and the client fetches presigned
part URLs instead.

| Method   | Path                             | Description                                      |
| -------- | -------------------------------- | ------------------------------------------------ |
| `OPTIONS`| `/api/v1/uploads`                | Advertise tus capabilities                       |
| `POST`   | `/api/v1/uploads`                | Open a session                                   |
| `HEAD`   | `/api/v1/uploads/{token}`        | Current offset                                   |
| `PATCH`  | `/api/v1/uploads/{token}`        | Send the next chunk (`tus` mode)                 |
| `POST`   | `/api/v1/uploads/{token}/parts`  | Mint presigned URLs for numbered parts (`direct`) |
| `POST`   | `/api/v1/uploads/{token}/complete` | Finish a `direct` upload, optionally with the browser's checksum claim |
| `DELETE` | `/api/v1/uploads/{token}`        | Abort an upload                                  |

## Content origin

These are served from the content host (or the app host in single-origin mode),
without a session cookie. File downloads (`/d/{token}/{name}`) support `Range`;
the streamed zip (`/z/{token}/{name}`) does not, because an archive built on the
fly has no length to declare and cannot be ranged into. The trailing filename is
decoration: only the token is authoritative. On object storage, a file download
is a redirect to a presigned bucket URL rather than a stream.

| Method             | Path                    | Description                                              |
| ------------------ | ----------------------- | -------------------------------------------------------- |
| `GET`/`HEAD`       | `/d/{token}/{name}`     | Signed download                                          |
| `GET`/`HEAD`       | `/z/{token}/{name}`     | Streamed zip of a selection or folder                    |
| `GET`/`HEAD`       | `/s/{token}`            | Public share link; a folder share takes `?p=` for a subpath |
| `POST`             | `/s/{token}`            | Submit the password of a protected share                 |

Older share links in the form `/s/{token}/{name}` still resolve.
