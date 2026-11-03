# Audit log

ByteWing records every sensitive action with a timestamp and the account that
performed it, so an administrator can answer what happened and when.

## What is recorded

- Sign-ins that succeed, by password or passkey, and sign-outs
- Password changes, two-factor being turned on or off, and passkey
  registration or removal
- Invitations created or revoked, accounts joined, changed or deleted
- Access grants and revocations, group creation, deletion and membership
- Share links created and revoked
- Files moved to the trash, restored, purged, or the trash emptied
- API token creation and revocation
- A change to the retention policy

A failed sign-in is not written to the log: there is no account to attribute it
to. Failed attempts are rate-limited per address and per account instead, and a
failed attempt against a real account is indistinguishable from one against a
name that does not exist.

## Reading it

The audit log is visible to administrators from the **Activity** screen, with
keyset pagination for a large history. Each entry carries the actor, the action,
the target and the client IP.

Because the IP is recorded, the log is data you should not keep indefinitely.
Under **Settings → Retention** an administrator can set a number of days after
which old entries (and expired trash) are purged automatically.
