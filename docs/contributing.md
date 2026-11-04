# Contributing

ByteWing is Apache-2.0 and welcomes contributions. The full guide lives in the
repository at
[`CONTRIBUTING.md`](https://github.com/andersmmg/bytewing/blob/main/CONTRIBUTING.md);
this page is a short orientation.

## Working on the code

- Backend is Go 1.27 with pure-Go SQLite; run `make check` before a PR.
- Frontend is React + Vite in `web/`.
- Commits follow [Conventional Commits](https://www.conventionalcommits.org)
  (`feat:`, `fix:`, `docs:`, …); releases are cut automatically by
  release-please from those messages.

## Working on these docs

This site is built with [VitePress](https://vitepress.dev). From `docs/`:

```sh
pnpm install
pnpm dev        # local preview with hot reload
pnpm build      # production build; fails on a dead internal link
```
