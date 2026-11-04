# Installation

ByteWing ships as a single container image on GitHub Container Registry. The
quickest way to run it is Docker Compose.

## Requirements

- Docker with the Compose plugin
- A directory on the host for your files
- Two hostnames pointing at the server (recommended; see
  [why two origins](#one-origin-or-two))

## Docker Compose

Create a `docker-compose.yml`:

```yaml
services:
  bytewing:
    image: ghcr.io/andersmmg/bytewing:latest
    restart: unless-stopped
    environment:
      BYTEWING_APP_URL: https://bytewing.example.tld
      BYTEWING_CONTENT_URL: https://dl.example.tld
      BYTEWING_DATA_DIR: /data
      BYTEWING_CONFIG_DIR: /config

    volumes:
      # A directory on the host, not a Docker volume: that is what gives ByteWing
      # the whole disk instead of a partition of it. ByteWing adopts whoever owns
      # this directory — no user ids to set. (It won't run as root: chown a
      # root-owned folder to any non-root user first.)
      - /mnt/storage:/data
      # The database lives apart from the browsable tree.
      - bytewing-config:/config

    ports:
      - "8080:8080"

volumes:
  bytewing-config:
```

Start it:

```sh
docker compose up -d
docker compose logs bytewing
```

::: tip Pin a version
`:latest` follows every release. For a reproducible deployment, pin a tag such
as `ghcr.io/andersmmg/bytewing:0.11.0`.
:::

## First run

**There is no default account.** Shipping known credentials would leave every
instance exposed between deployment and the moment someone changes them.

Instead, the log prints a **one-time setup link** on first start:

```
open http://localhost:8080/setup?token=… to create the administrator
```

Open it, create the administrator account, and you are in. The link expires
after 24 hours; restart the container to mint a fresh one if you miss it.

## One origin, or two?

ByteWing serves two things: the **application** (the UI and API, which use a
session cookie) and the **content** (the files themselves). Serving them from
two different hostnames is what guarantees a downloaded or shared file can never
carry a session cookie.

- **Two origins (recommended).** Set both `BYTEWING_APP_URL` and
  `BYTEWING_CONTENT_URL` to different hostnames. Enables in-place preview of
  images, video, audio, text and PDF.
- **One origin (supported).** Leave `BYTEWING_CONTENT_URL` unset. The instance
  hardens itself automatically but sends every file as a download instead of
  previewing it in place. On the object-storage backend this does not apply:
  files are answered by presigned redirects to the store's own origin, so
  previews work either way.

Files can also live in an S3-compatible bucket instead of `/data`; see
[Object storage](/guide/object-storage).

See [Deployment](/guide/deployment) for a reverse-proxy example, and
[Environment variables](/reference/environment) for the full list.

## Building from source

To run exactly what is on a branch rather than a published image, the repository
ships a compose file that builds the image locally:

```sh
docker compose -f deploy/docker-compose.build.yml up -d
```

Or build the binary directly (requires Go 1.27 and pnpm):

```sh
make dist   # builds the web UI and embeds it, then compiles the binary
./bin/bytewing
```
