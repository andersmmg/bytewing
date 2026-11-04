# Permissions

Access is granted **per path**, to a user or a group, as a set of rights:
**read**, **write**, **delete**, **share** and **manage**. Read, write and
delete are checked in the storage layer on every filesystem operation, not in
an HTTP handler, so a forgotten endpoint cannot open a hole. **Share** and
**manage** are capabilities rather than filesystem operations, so the storage
guard does not check them: `Engine.Allows` performs the share check before a
link is created, and only administrators change rules.

## How a decision is made

Permissions are resolved **per bit**, not per rule. For each of the five rights,
ByteWing considers the rules that apply to the path (the user's own and those of
their groups), and the deepest rule that mentions that bit decides it. A denial
wins a tie at the same depth, so two rules that disagree about a path never
silently grant access.

This is what lets a rule granting read and write at `/projects` and one denying
write at `/projects/archive` coexist: the subject can read and write throughout
`/projects`, and the subtree stays read-only, without either rule restating the
other.

**Ownership** is layered on top: the owner of a file implicitly holds read,
write, delete and share on it, but not manage; uploading a file is no reason to
control who else may see it. A rule at any depth still overrides ownership, so
an administrator can lock a path even against the person who uploaded it.

## Traversal

A folder you were granted is reachable through the directories above it, without
those parents exposing their other contents. You see the path down to your grant
and nothing else.

## Inspecting access

The interface shows the rules set on a path, with each rule's subject,
permissions, whether it is recursive, and whether it is a denial. The API also
answers what the **caller** effectively holds at a path
(`GET /api/v1/permissions?path=`), which is what enables or disables actions in
the interface.

Because inheritance is per bit, the answer is a list of rules rather than a
single reason: trace the right you care about to the deepest rule that mentions
it, and a denial at that depth wins.

## Two rules that hold everywhere

- **The storage layer checks filesystem access**, not HTTP handlers.
- **No escalation by sharing**: you cannot share more than you hold. A
  read-only user cannot create a writable link.
