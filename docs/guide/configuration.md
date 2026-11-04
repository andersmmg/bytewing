# Configuration

Everything is configured through **environment variables**. A self-hosted
service is deployed from a compose file far more often than from a config file,
and one source of truth avoids the question of which wins.

Validation is strict and happens at startup: a misconfigured instance refuses to
start with a clear message rather than failing later in a way that looks like a
bug.

## Minimal configuration

The smallest valid configuration sets the public address and where data and the
database live:

```yaml
environment:
  BYTEWING_APP_URL: https://bytewing.example.tld
  BYTEWING_DATA_DIR: /data
  BYTEWING_CONFIG_DIR: /config
```

`BYTEWING_CONFIG_DIR` must **not** be inside `BYTEWING_DATA_DIR`: a database users
can browse is a database they can download, password hashes included. ByteWing
refuses to start otherwise.

## Key choices

- **Two origins vs one**: set `BYTEWING_CONTENT_URL` to a second hostname to
  enable in-place previews. See [One origin, or two?](/guide/installation#one-origin-or-two).
- **HTTPS**: session cookies are marked `Secure` automatically when
  `BYTEWING_APP_URL` uses `https`. A `Secure` cookie sent over plain HTTP is
  silently dropped by the browser, so always terminate TLS in front of ByteWing
  in production.
- **File ownership**: nothing to configure. ByteWing adopts whoever owns the data
  directory you mount, so a folder you own just works. It will not run as root:
  if the mounted folder is root-owned, `chown` it to any non-root user once.
- **Storage backend**: local disk by default. `BYTEWING_STORAGE_BACKEND: s3`
  keeps files in an S3-compatible bucket instead; then `BYTEWING_DATA_DIR` is not
  read and `BYTEWING_S3_BUCKET` + `BYTEWING_S3_REGION` are required. The full
  setup (providers, credentials, CORS and what changes) is under
  [Object storage](/guide/object-storage).

## Object storage mode

In S3 mode there is no local disk to fill, so `BYTEWING_RESERVE_BYTES` does not
apply to file writes, and `BYTEWING_READ_ONLY` plus the instance
[quota](/reference/environment#bytewing-quota-bytes) are the write guards. Two
things still touch local disk: the SQLite database in `BYTEWING_CONFIG_DIR`, and
the staging directory for URL downloads.

| Operation | Local | S3 |
|---|---|---|
| Browse, upload, download, trash, move, copy | ✓ | ✓ |
| In-place previews | needs two origins | ✓, either mode |
| Text preview | ✓ | ✓, bounded read |
| Thumbnails | ✓ | ✓, background, oversized sources skipped |
| Fetch URL to server | ✓ | ✓, via local staging |
| Extract a ZIP | ✓ | refused |
| Checksum (SHA-256) | ✓ | ✓, from the browser at upload time or the store's own value; computed once in the background otherwise |

The full list, with defaults and validation rules, is in the
[Environment variables reference](/reference/environment), and the behaviour in
detail is under [Object storage](/guide/object-storage).
