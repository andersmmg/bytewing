# Deployment

ByteWing is one container behind a reverse proxy that terminates TLS. The only
ByteWing-specific considerations are the two origins and, with object storage,
the bucket's CORS rule.

## Reverse proxy

Point both hostnames at the same container port (`8080`). ByteWing decides which
role a request plays from its `Host` header, so no path rules are needed.

```
bytewing.example.tld  ─┐
                     ├─►  bytewing:8080
dl.example.tld      ─┘
```

Any proxy works (Caddy, Traefik, nginx). Terminate TLS at the proxy and forward
plain HTTP to the container. Because ByteWing derives `Secure` cookies from
`BYTEWING_APP_URL`, make sure that variable uses `https`.

## Dokploy

ByteWing deploys cleanly on [Dokploy](https://dokploy.com). Use the reference
compose file, set the two domains in the Dokploy UI, and point the storage
volume at a host directory. Dokploy handles TLS and routing for both hostnames.

## Forwarded headers

ByteWing trusts forwarded headers only when the request's direct peer is a
private or loopback address: in practice, a reverse proxy on the same host or
network. Then it prefers `X-Real-Ip`, and otherwise takes the **rightmost**
entry of `X-Forwarded-For`: the address the proxy itself saw connect, which a
client cannot forge past it. Any values further left are ignored.

A request that arrives from a public address has its headers ignored, so
putting the instance directly on the internet does not let anyone spoof their
rate limit or audit entry. Behind a proxy, make sure it appends the client
address rather than passing a client-supplied header through unchanged.

## Object storage

For S3, Backblaze B2, R2 and MinIO, see
[Object storage](/guide/object-storage) for configuration, credentials, the
required bucket CORS rule and the operations that behave differently. The short
version: the browser talks to the bucket directly through presigned URLs, so the
bucket needs a CORS rule for the application origin, and it can stay private.

## Backups

Back up two things separately: your **data directory** (the files; use your usual
disk backup) and the **database**, which holds every account, permission,
share and token. Losing the database loses all of that even though the files
survive, so it is the critical half.

ByteWing has built-in commands for the database. `backup` is safe to run while the
server is live: it takes a consistent snapshot with SQLite's `VACUUM INTO`:

```sh
docker compose exec bytewing bytewing backup
# → /config/backups/bytewing-2026-08-03-020000.db
```

Pass a path to choose where it lands, which is handy for a nightly cron:

```sh
docker compose exec bytewing bytewing backup /config/backups/nightly.db
```

The snapshot is a self-contained SQLite file you can copy off the host.

To restore, **stop the server first**. Replacing the database under a running
instance would corrupt it. `restore` validates the snapshot and copies the
current database aside (as `bytewing.db.pre-restore-…`) before replacing it:

```sh
docker compose stop bytewing
docker compose run --rm bytewing bytewing restore /config/backups/nightly.db
docker compose start bytewing
```

If the restored database references files that no longer exist on disk, because
the snapshot predates changes to the tree, `restore` lists them so you know what
drifted. Schema migrations also snapshot the database automatically before
they run, so an upgrade always has a point to roll back to.
