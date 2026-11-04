# Environment variables

Every setting comes from the environment and is validated at startup.

## Application

### `BYTEWING_APP_URL`

**Required.** The public address of the application origin (the UI and API),
e.g. `https://bytewing.example.tld`. Used to build links and to decide whether
session cookies may be marked `Secure`: cookies are `Secure` when this URL uses
`https`.

### `BYTEWING_CONTENT_URL`

_Optional._ The public address of the content origin, on a **different host**
from `BYTEWING_APP_URL`. When set, files are served from this cookieless origin
and can be previewed in place. When unset, the instance runs in single-origin
mode: it hardens itself automatically but sends every file as a download.

The hardening only applies to backends that stream files through the
application itself. With the object-storage backend, downloads are presigned
redirects to the store's own origin, so inline preview works regardless of
this variable.

ByteWing refuses to start if this shares a host with `BYTEWING_APP_URL`: separate
origins are what keep an uploaded file from reaching a session.

### `BYTEWING_ALLOW_URL_FETCH`

_Optional, default `false`._ Lets the server download files from user-supplied
URLs using the interface's **Download from URL** import.

When on, the download runs as a background job with the usual defences:
public addresses only, revalidated on every redirect hop, and is bounded by
free disk space on both the local download staging volume and the storage space
quota.

### `BYTEWING_MIN_PASSWORD_LENGTH`

_Optional, default `10`._ The shortest password an account may choose. The
value must be a whole number from `3` to `1024`. The setup, invite,
password-reset and change-password forms apply the same limit.

## Storage

### `BYTEWING_DATA_DIR`

**Required unless the storage backend is `s3`.** The storage root, meaning the
directory users browse. Point it at a host directory (a bind mount), not a
Docker volume, to give ByteWing the whole disk. It is not read in object-storage
mode.

### `BYTEWING_CONFIG_DIR`

**Required.** Where the SQLite database lives. It must **not** be inside
`BYTEWING_DATA_DIR`; ByteWing enforces this, because a database users can list is a
database they can download, password hashes included.

### `BYTEWING_RESERVE_BYTES`

_Optional._ Free space kept in hand before refusing writes. Below the threshold
the service goes read-only (reads, sign-in and deletes still work) rather than
letting the disk fill until SQLite can no longer write. Defaults to a built-in
reserve. With the object-storage backend this setting is unused: there is no
local disk to fill, and the staging area is bounded by its own free space.

### `BYTEWING_READ_ONLY`

_Optional, default `false`._ When `true`, refuses every write regardless of free
space. Useful for maintenance or a read-only mirror.

### `BYTEWING_PREVIEW_MAX_BYTES`

_Optional, default `2MiB`._ Bounds server-side preview reads: text previews
are truncated past it, and thumbnail sources past it are refused in
object-storage mode (where every preview byte is fetched remotely). Accepts a
bare byte count or a `KiB`/`MiB`/`GiB` suffix.

### `BYTEWING_THUMB_MAX_BYTES`

_Optional, default `512MiB`._ Caps the source file a background thumbnail job
will fetch. Thumbnails themselves are small JPEGs kept under a reserved
prefix in the same storage; larger files keep their icon. Same size format.

### `BYTEWING_QUOTA_BYTES`

_Optional, default `0` (unlimited)._ Caps the bytes of live files the instance
will hold: uploads that would push usage past it are refused with a "storage
quota exceeded" error, checked when the upload starts and again when it
publishes. Per-account limits are set by an admin under **Members** (zero
means unlimited) and are enforced the same way. Usage counts live files only
(trash and generated thumbnails do not count), and the sidebar shows each
account its own usage against its own limit. Same size format.

::: tip Retention
How long the audit log and the trash are kept is **not** an environment setting:
it is configured at runtime under **Settings → Retention** (admin). See the
[audit log](/features/audit-log).
:::

## Object storage

_Optional._ Set `BYTEWING_STORAGE_BACKEND=s3` to keep files in S3-compatible
object storage (AWS S3, Backblaze B2, Cloudflare R2, MinIO) instead of
`BYTEWING_DATA_DIR`. `BYTEWING_S3_BUCKET` and `BYTEWING_S3_REGION` are required in
this mode; `BYTEWING_S3_ENDPOINT` addresses a non-AWS provider, and
`BYTEWING_S3_ACCESS_KEY_ID` / `BYTEWING_S3_SECRET_ACCESS_KEY` carry static
credentials (unset means the ambient AWS chain). `BYTEWING_S3_PREFIX` scopes the
tree inside the bucket, `BYTEWING_S3_FORCE_PATH_STYLE` switches to path-style
addressing, `BYTEWING_S3_PUBLIC_ENDPOINT` is the bucket's address as clients see
it (for presigned URLs), `BYTEWING_S3_PART_SIZE` is the multipart part length
(default 64 MiB, min 5 MiB; the browser uploads several parts concurrently), and
`BYTEWING_S3_PRESIGN_TTL` bounds minted capability URLs (default 15 min).

Setup, provider examples, bucket CORS and the behaviour that changes are covered
in [Object storage](/guide/object-storage).

### `BYTEWING_S3_META_TTL`

_Optional, default `1m`._ How long directory listings and file stats are served
from the metadata cache before one S3 request revalidates them. Changes made
through the app invalidate immediately; the TTL only bounds how stale rows for
changes made outside the app can get. Accepts a Go duration (`30s`, `5m`);
`0` disables expiry-based revalidation.

## Server

### `BYTEWING_LISTEN`

_Optional, default `:8080`._ The address the HTTP server binds.

## File ownership

There is nothing to configure, and no `PUID`/`PGID`. The container entrypoint
runs the server as **whoever owns the mounted data directory**, so a bind-mounted
folder is written by its own owner and the first upload never fails on a
permission error. Mount a directory you own and it works.

Two edge cases the entrypoint handles:

- **A fresh, empty mount** (Docker creates a missing bind-mount target owned by
  root) is claimed for the image's built-in unprivileged user automatically.
- **A root-owned directory that already holds files** is not reowned. ByteWing
  will not run as root. `chown` it to any non-root user once
  (`chown -R 1000:1000 <host directory>`) and ByteWing follows.

With the object-storage backend there is no data directory to adopt, so the
container runs as its built-in unprivileged user (uid/gid `1000`) and the config
volume is chowned to match. Mount `/config` and nothing else.

To pin a specific user regardless, set `user:` on the service in your compose
file; the entrypoint sees it is already unprivileged and adopts nothing.
