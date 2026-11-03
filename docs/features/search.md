# Search

Search finds entries by name across a folder and everything below it, showing
only what you are allowed to see.

## How it works

There is no index to go stale. Each search walks the live tree from where you
start it, so a file dropped in over SSH is findable immediately, and a result
can never describe something the disk no longer has.

Every entry the walk meets passes through the same per-path authorisation as a
listing, so a match in a folder you were granted is returned without exposing
anything else in the directories above it.

## What it matches

Names only, not file contents, matched as a case-insensitive substring, so
`report` finds `Q3-report.pdf` and `REPORT.md`.

Paths are stored in composed Unicode form (NFC). An accented name matches when
the query is typed that way; a decomposed form with separate combining marks
will not match, so type the character rather than pasting a decomposed
version.

The walk is breadth-first, which means matches nearer your starting folder come
first rather than after an alphabetical sweep of the whole tree.

## Limits

A search is bounded so one request cannot scan a pathological tree forever:

- `limit` defaults to 200 results and is capped at 500.
- The walk descends into at most 20,000 directories.

When either limit is reached the response carries `truncated: true` and the
results found so far. Directories sort before files, then by name. A directory
that disappears or becomes unreadable mid-walk is skipped rather than failing the
whole search.

Narrowing the starting `path`, or a more specific query, is how you get complete
results from a huge tree.
