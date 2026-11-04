# Object storage (S3)

ByteWing can keep your files in an S3-compatible bucket instead of a directory on
the host. The application still runs as one binary; only where the bytes live
changes.

This works with AWS S3, Backblaze B2, Cloudflare R2, MinIO and anything else that
speaks the S3 API.

## What changes

- **Files live in the bucket.** `BYTEWING_DATA_DIR` is not read in this mode.
- **The database stays on local disk.** Accounts, permissions, shares, the audit
  log and the metadata cache are still SQLite files under
  `BYTEWING_CONFIG_DIR`.
- **The browser transfers bytes directly with the bucket.** Uploads are
  presigned multipart `PUT`s; downloads and previews are a redirect to a
  presigned `GET`. The application mints the URLs and never proxies the bytes.
- **Presigned URLs are temporary.** They are valid for
  `BYTEWING_S3_PRESIGN_TTL` (15 minutes by default). A share link re-mints one
  on every request, so the link itself does not expire with the URL.
- **Zip downloads are the exception.** A folder or multi-file zip has to be
  assembled on the fly, so the server reads the objects and streams the archive.

## Configuration

```yaml
services:
  bytewing:
    image: ghcr.io/andersmmg/bytewing:latest
    restart: unless-stopped
    environment:
      BYTEWING_APP_URL: https://bytewing.example.tld
      BYTEWING_CONFIG_DIR: /config
      BYTEWING_STORAGE_BACKEND: s3
      BYTEWING_S3_BUCKET: bytewing
      BYTEWING_S3_REGION: us-west-004
      BYTEWING_S3_ENDPOINT: https://s3.us-west-004.backblazeb2.com
      BYTEWING_S3_ACCESS_KEY_ID: …
      BYTEWING_S3_SECRET_ACCESS_KEY: …

    volumes:
      # The database and the staging area. No data directory is mounted.
      - bytewing-config:/config

volumes:
  bytewing-config:
```

| Variable | Default | Meaning |
| --- | --- | --- |
| `BYTEWING_STORAGE_BACKEND` | `local` | Set to `s3` to use a bucket. |
| `BYTEWING_S3_BUCKET` | — | **Required.** The bucket that holds the tree. |
| `BYTEWING_S3_REGION` | — | **Required.** The bucket's region. On R2 this is `auto`. |
| `BYTEWING_S3_ENDPOINT` | AWS | The provider's S3 address, as the server reaches it. Leave unset for AWS. |
| `BYTEWING_S3_PUBLIC_ENDPOINT` | `BYTEWING_S3_ENDPOINT` | The same address as **clients** reach it, used only for presigned URLs. Set it when the server uses an internal address (a MinIO service name, say) that a browser cannot resolve. |
| `BYTEWING_S3_PREFIX` | none | Scopes the tree to a key prefix, e.g. `bytewing/`. Leading and trailing slashes are trimmed. |
| `BYTEWING_S3_ACCESS_KEY_ID` | ambient chain | Static credential. Set it together with the secret; leave both unset to use the AWS chain (instance role, environment, shared config). |
| `BYTEWING_S3_SECRET_ACCESS_KEY` | ambient chain | The secret half of the pair above. |
| `BYTEWING_S3_FORCE_PATH_STYLE` | `false` | Addresses buckets as `endpoint/bucket/key` instead of `bucket.endpoint`. MinIO and some gateways need `true`. |
| `BYTEWING_S3_PART_SIZE` | `64MiB` | Length of each multipart part except the last. Minimum `5MiB`, which is a provider-side rule. |
| `BYTEWING_S3_PRESIGN_TTL` | `15m` | How long a minted download or upload URL stays valid. |
| `BYTEWING_S3_META_TTL` | `1m` | How long cached listings and stats are served before one S3 request revalidates them. `0` disables expiry-based revalidation. |

The bucket does **not** need to be public. Presigned URLs work against a private
bucket, and a public one would let anyone who knows a key skip every permission
check. Leave public access off.

## Provider examples

| Provider | `ENDPOINT` | `REGION` | Notes |
| --- | --- | --- | --- |
| Backblaze B2 | `https://s3.<region>.backblazeb2.com` | the bucket's region, e.g. `us-west-004` | Use an application key scoped to the bucket. |
| AWS S3 | unset | the bucket's region | An instance role is enough; leave the static keys unset. |
| Cloudflare R2 | `https://<account-id>.r2.cloudflarestorage.com` | `auto` | Create an S3 API token for the bucket. |
| MinIO | your server's S3 port | as configured | Usually `BYTEWING_S3_FORCE_PATH_STYLE: true`. |

### Credentials

A static key needs the usual object permissions on that one bucket:

- `s3:ListBucket`: directory listings and existence checks
- `s3:GetObject`: downloads, previews, thumbnails and checksums
- `s3:PutObject`: uploads, copies, thumbnails
- `s3:DeleteObject`: deletes, purges and the delete half of a move
- `s3:AbortMultipartUpload`, `s3:ListMultipartUploadParts`: resumable uploads

The server checks the bucket at startup (`HeadBucket`) and refuses to start if it
cannot reach it, which is preferable to a server that runs and fails on the
first upload.

## Browser access and CORS

Because uploads and previews come from the browser, the bucket needs a CORS rule
covering the **application origin**:

- `GET`, `PUT`, `HEAD`
- any request header: part uploads send `Content-Type`, previews send `Range`
- `ETag`, `Content-Length`, `Content-Range` and `Accept-Ranges` exposed

The repository ships the rule in Backblaze's native JSON format as
`deploy/b2-cors.json`:

```json
[
  {
    "corsRuleName": "bytewing-dev",
    "allowedHeaders": ["*"],
    "allowedOperations": ["s3_get", "s3_put", "s3_head"],
    "allowedOrigins": ["http://localhost:8080"],
    "exposeHeaders": ["ETag", "Content-Length", "Content-Range", "Accept-Ranges"],
    "maxAgeSeconds": 3600
  }
]
```

The file targets local development (`http://localhost:8080`). Replace the origin
with your application origin (`https://bytewing.example.tld` here) before
applying it. With Backblaze's CLI (the binary is `bbb2` on systems where `b2` is
already a build tool):

```sh
bbb2 account authorize <applicationKeyId> <applicationKey>

# bucketType is positional and required; repeat the bucket's current type
# (allPrivate) so the update preserves it.
bbb2 bucket update --cors-rules "$(<deploy/b2-cors.json)" <bucketName> allPrivate
```

On AWS, R2 and MinIO the same rule goes in the bucket's CORS settings in the
console or with `aws s3api put-bucket-cors`.

If uploads fail with a CORS error and downloads work, the rule is missing `PUT`
or has the wrong origin. If previews load but seeking fails, `Range` is not in
`allowedHeaders` or `Content-Range` is not exposed.

## What works in this mode

| Operation | Object storage |
| --- | --- |
| Browse, upload, download, trash, move, copy | ✓ |
| In-place previews | ✓, even in single-origin mode |
| Text preview | ✓, bounded read |
| Thumbnails | ✓, generated in the background; oversized sources are skipped |
| Download from URL | ✓, needs local staging space (below) |
| Extract a ZIP | ✗ (refused) |
| SHA-256 checksum | ✓, from the browser at upload time, the store's own value, or a background job |
| Instance quota, read-only mode | ✓ |
| Free-space reserve | not applicable; there is no local disk to fill |

Extracting an archive is not available. A zip stored in the bucket can be
downloaded and extracted locally, but not unpacked in place.

Browsing is served from a metadata cache. Listings and stats are mirrored into
SQLite, so a folder you have already visited is one query rather than a fresh
round trip. Changes made through ByteWing invalidate the cache immediately;
objects added or removed outside it appear after `BYTEWING_S3_META_TTL`.

## Local staging

Some work still needs a local file: a URL download is fetched to disk and then
uploaded, and a non-direct upload passes through the server. That happens under
`$BYTEWING_CONFIG_DIR/stages`.

Size that volume for the largest file you expect to fetch from a URL, not for
your whole library. The fetch is refused when free space on that volume is
exhausted. A partial upload stops being resumable after its session expires, but
the staged bytes are not swept automatically today: if the directory grows,
clear `stages` while the server is stopped.

## Backups

Back up both halves separately:

1. **The bucket**: use your provider's versioning or its own backup tool. Object
   storage is not a backup by itself; a deleted prefix is gone.
2. **The database**: `bytewing backup` writes a consistent SQLite snapshot and
   is safe to run while the server is live. See
   [Backups](/guide/deployment#backups).

Set a lifecycle rule aborting incomplete multipart uploads (7 days is sane).
Aborted or abandoned uploads leave parts behind that you pay for but cannot see
in a listing.

## Moving existing files

There is no built-in local-to-bucket migration command. Copy the tree into the
bucket (or its configured prefix) preserving relative paths, then switch
`BYTEWING_STORAGE_BACKEND` and restart. Permissions, accounts and shares live in
the database and keep applying, because they are keyed by path.

If the paths do not line up, the files will not appear where their rules expect
them. Uploading through the interface instead is slower but avoids that class of
mistake.

Checksums are not carried over by a manual copy: the cached digests describe the
files the database saw. Request a checksum again after the move and the
background job computes a fresh one.

## Troubleshooting

- **`403 Forbidden` from the bucket on startup**: the key lacks `ListBucket`, or
  the bucket name, region or endpoint is wrong.
- **`SignatureDoesNotMatch`**: the region does not match the bucket's, or the
  key pair was pasted with whitespace.
- **Uploads hang or fail at the first part**: the CORS rule is missing `PUT`,
  does not list the app origin, or the endpoint is not reachable from the
  browser. Check `BYTEWING_S3_PUBLIC_ENDPOINT`.
- **Previews fail to load but files download**: the presigned URL points at an
  address only the server can resolve. Set `BYTEWING_S3_PUBLIC_ENDPOINT`.
- **New objects do not appear**: they were added outside ByteWing. Wait out
  `BYTEWING_S3_META_TTL`, or set it to `0` if external writes are routine.
