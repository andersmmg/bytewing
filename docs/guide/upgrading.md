# Upgrading

ByteWing follows [semantic versioning](https://semver.org). While it is below
`1.0`, a minor bump (`0.5 → 0.6`) may contain new features and, occasionally, a
breaking change called out in the release notes.

## With Docker

Pull the new image and recreate the container:

```sh
docker compose pull
docker compose up -d
```

If you pinned a version tag, bump it first:

```yaml
image: ghcr.io/andersmmg/bytewing:v0.10.0
```

## Database migrations

Schema migrations run automatically at startup. ByteWing takes an automatic
snapshot of the database before applying one, so an interrupted upgrade does not
leave a half-migrated database.

::: warning Back up first
Migrations are one-way. Before a major upgrade, make sure you have a copy of the
config volume; see [Backups](/guide/deployment#backups).
:::

## Checking the running version

The version is shown in the sidebar of the interface, and returned by the
config endpoint (which requires a session or API token):

```sh
curl -H "Authorization: Bearer $TOKEN" \
  https://bytewing.example.tld/api/v1/config
```
