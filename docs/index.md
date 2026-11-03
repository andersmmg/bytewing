---
layout: home

hero:
  name: ByteWing
  text: The file server that should have existed
  tagline: A self-hosted file server in one binary, with resumable transfers, permissions enforced in the storage layer and sessions you can revoke. Apache-2.0.
  actions:
    - theme: brand
      text: Get started
      link: /guide/installation
    - theme: alt
      text: Introduction
      link: /guide/introduction
    - theme: alt
      text: View on GitHub
      link: https://github.com/andersmmg/bytewing

features:
  - title: One binary
    details: Embedded web UI, pure-Go SQLite, no required runtime dependencies and no database server. Deploy one container.
  - title: Your files on your disk
    details: Point ByteWing at a directory or bind mount, with no extra service. Or use an S3-compatible bucket instead.
  - title: Fast, resumable transfers
    details: Resumable uploads (tus) and streamed zip downloads, limited by disk and network rather than CPU.
  - title: Real permissions
    details: Granular access per path, granted to a user or a group, enforced in the storage layer.
  - title: Sessions that end
    details: Opaque server-side tokens, never JWT. Logging out revokes access immediately.
  - title: Share links
    details: Revocable links with expiry and an optional password, served cookieless from a separate origin with full Range support.
  - title: Trash and an audit trail
    details: Deleted files wait in a trash you can restore from, and an audit log records who changed what. How long both are kept is configurable.
  - title: No telemetry
    details: No usage metrics, no update checks, no outbound request you didn't ask for.
---
