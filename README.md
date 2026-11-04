<p align="center">
  <img src="docs/public/logo.png" alt="ByteWing" width="128">
</p>

<h1 align="center">ByteWing</h1>

A self-hosted file server in a single binary. Store, organise and share large files
on your own server or an S3-compatible bucket, with a permission model and a JSON API
included.

---

> [!NOTE]
> **ByteWing is at 0.x, and a live instance serves the author's files in
> production.** The data model is stable; expect the occasional rough edge and
> watch the releases for breaking changes, which land only on a minor bump while
> 0.x.

---

## Why

[File Browser](https://github.com/filebrowser/filebrowser) was archived on 1 September 2026.
Its [author explained why](https://hacdias.com/2026/07/28/filebrowser/) it could not be fixed by
patching: written without a security model or API design, it would need a rewrite. Known
vulnerabilities remain public and unfixed.

ByteWing is that rewrite.

## What it does

**Files**

- Browse as a list or a grid, sorted and grouped, with thumbnails and a
  recursive search that only shows what you are allowed to see.
- **Resumable uploads** (tus), because a 40 GB transfer that fails at 90% must not
  start over. Drag in files or a whole folder tree.
- Copy, move, rename, and a **trash** you can restore from. Copying a folder or a
  very large file runs as a background job with a progress bar.
- Download through short-lived signed links that work with download managers.
- **Checksums** on demand: a file's SHA-256, computed as a background job and
  cached, so verifying even a huge file never blocks the server.
- **Extract a ZIP** into a new folder, as a background job. A booby-trapped
  archive cannot escape: entry names go through the same path validation as
  everything else, symlinks are ignored, and a decompression bomb is refused by
  caps on entry count, total size and per-entry expansion.
- **Download from a URL** straight to the server, at its bandwidth rather than
  yours, which matters when the file is tens of gigabytes. Off by default;
  `BYTEWING_ALLOW_URL_FETCH` turns it on. The fetch is guarded against
  server-side request forgery: the real IP being dialled is validated on every
  connection and every redirect, so a URL that resolves or redirects to an
  internal address is refused.

**Sharing**

- Unique, revocable **share links** for a file or a folder: an optional expiry,
  and (for a single file) an optional password.
- Links work with download managers. No JavaScript, no cookies, correct `Range`
  support, and parallel connections. Served from a separate origin so a shared file
  can never reach a session.

**Multiple people**

- Invite users with a one-time link; they create their own account.
- **Granular permissions per path**, granted to a user or a **group**. A folder
  granted to a group reaches everyone in it, and a folder you were granted is
  reachable through the directories above it without exposing their contents.
- Admins manage accounts (promote, disable, remove) and groups; everyone manages
  their own password and sessions, shown by device and **location** rather than a
  raw IP.
- **Two-factor authentication** (TOTP) and single-use **recovery codes** to reset a
  forgotten password, with no email anywhere.
- **Passkeys** (WebAuthn) alongside the password: register one in settings, then
  sign in with a biometric or PIN instead of typing anything. The password keeps
  working, and managing passkeys always asks for the password again.
- An **audit log** records every sensitive action; how long the log and the trash
  are kept is configurable.

**Throughout**

- **One binary.** Embedded web UI, pure-Go SQLite, no database server. Video
  thumbnails use `ffmpeg` when it is installed; nothing else needs it.
- **Sessions that end.** Opaque server-side tokens, never JWT. Logging
  out revokes the token immediately.
- **API tokens** for scripts, CI and integrations. A token carries your own
  permissions, is shown once, and can be revoked at any time.
- **Back up and restore** the database from one command, with an automatic
  snapshot taken before every schema migration.
- **No telemetry.** There are no usage metrics or update checks, and no outbound
  request you didn’t ask for.

## What it is not

- Not a collaborative document editor (no OnlyOffice, no Collabora)
- Not a groupware suite: no calendar, contacts, notes or chat
- No plugin marketplace
- No bidirectional desktop sync client

ByteWing aims to be the File Browser that should have existed, not a Nextcloud alternative.

## Stack

Go 1.27 · SQLite ([`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite), pure Go) ·
React 19 + Vite · [shadcn/ui](https://ui.shadcn.com/) (Radix + Tailwind CSS 4) ·
[Phosphor Icons](https://phosphoricons.com/)

## Documentation

Full documentation, covering installation, configuration, features and the API,
lives at **[andersmmg.github.io/bytewing](https://andersmmg.github.io/bytewing/)**.
The source is in [`docs/`](docs/) ([VitePress](https://vitepress.dev)); build it
locally with `cd docs && pnpm install && pnpm dev`.

## Deploying

ByteWing ships as a container image. The compose file builds it from source, which
is what a self-hosted instance usually wants:

```sh
docker compose -f deploy/docker-compose.build.yml up -d
docker compose -f deploy/docker-compose.build.yml logs bytewing
```

The log prints a one-time setup link. Open it to create the administrator.

**There is no default account.** Shipping known credentials would leave every
instance exposed between deployment and the moment someone remembers to change
them, and scanners know the defaults of every self-hosted project. The setup
link is single-use, expires in 24 hours, and is replaced on every start. An old
one printed in a previous log stops working.

Two things decide whether a deployment works:

| Variable | Why it matters |
| -------- | -------------- |
| `BYTEWING_CONFIG_DIR` | Must not sit inside `BYTEWING_DATA_DIR`. ByteWing refuses to start otherwise: a database users can list is one they can download, password hashes included.                                                                                                                                                                            |
| `BYTEWING_CONTENT_URL` | A second hostname. Without it, user content is served from the application origin and ByteWing hardens itself: every file becomes an attachment and nothing renders in place. On the object-storage backend this does not apply, because downloads are presigned redirects to the store's own origin, so previews keep working. |

There are **no `PUID`/`PGID` to set**. ByteWing adopts whoever owns the data
directory you mount, so a folder you own just works. The one exception: it will
not run as root, so if the mounted folder is owned by root, `chown` it to any
non-root user once and ByteWing follows.

## Building

```sh
make dist    # build the interface and a binary embedding it
make check   # what CI runs: fmt, vet, test, govulncheck
```

Requires Go 1.27+, and Node with pnpm for the interface. The version is read from
[`version.txt`](version.txt), the single source of truth kept by release-please, so
a local build reports the same version as a release. See `make help` for targets.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). In short: work however you like on your
branch, open a pull request whose **title** is a
[Conventional Commit](https://www.conventionalcommits.org/); it is squash-merged
and becomes the release note.

## Credits

ByteWing began as [zefile](https://github.com/krishna2206/zefile) by
Fitiavana Anhy Krishna, and carries its git history forward under the same
[Apache-2.0](LICENSE) licence. See [NOTICE](NOTICE).

IP geolocation for the sessions screen uses the free
[IP to City Lite database by DB-IP](https://db-ip.com/db/download/ip-to-city-lite)
([CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)), read offline, so no
lookup ever leaves the server.

## Licence

[Apache-2.0](LICENSE).
