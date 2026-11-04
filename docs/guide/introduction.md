# Introduction

ByteWing is a self-hosted file server: a single binary that serves a directory of
files over a clean web interface and a JSON API. It aims to be **the File
Browser that should have existed**, not a Nextcloud or Google Drive replacement.

## Why it exists

[File Browser](https://github.com/filebrowser/filebrowser) proved that a
single-binary file server is the right shape. But it carries a few defects its
own maintainer has acknowledged: JWT sessions that cannot be revoked, slow
downloads, and a permission model that is easy to bypass. ByteWing keeps the shape
and fixes what made those defects possible:

- **Sessions that end.** Opaque server-side tokens, checked in the database on
  every request. Logging out revokes access on the next request, so a token does
  not keep working after sign-out.
- **Permissions in the storage layer.** Access checks live below the HTTP
  handlers, so forgetting an endpoint cannot open a hole.
- **Transfers limited by the disk.** Resumable uploads and streamed, uncompressed
  zip downloads are bound by disk and network rather than CPU.

## What it is

- **One binary.** Embedded web UI, pure-Go SQLite, no runtime dependencies, no
  database server.
- **Multi-user.** Invite people with one-time links; grant access per path to a
  user or a whole group; an audit log records who did what.
- **Shareable.** Revocable links with expiry and optional password, served
  cookieless from a separate origin.
- **API-first.** Everything the interface does is a plain JSON API, reachable
  with an [API token](/reference/api-tokens).

## What it is not

- Not a collaborative document editor (no OnlyOffice, no Collabora)
- Not a groupware suite: no calendar, contacts, notes or chat
- No plugin marketplace
- No bidirectional desktop sync client

## Next steps

- [Install ByteWing](/guide/installation) with Docker Compose
- Review the [configuration](/guide/configuration) and
  [environment variables](/reference/environment)
- Read about [permissions](/features/permissions) and [sharing](/features/sharing)
