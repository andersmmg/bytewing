# Sharing

A share link gives access to a file or folder without an account. Links are
unique, revocable, and served cookieless from the content origin.

## Creating a link

Share any file or folder you can read. A link can carry:

- an **expiry**, after which it stops working;
- for a file, an optional **password**, hashed with Argon2id;
- for a folder, browsing confined to the shared subtree.

A folder share cannot be password-protected yet: the password flow is built
around a single response, and browsing a folder needs a session-like grant that
is not there yet. You cannot share more than you hold:
[permissions](/features/permissions) are enforced at creation.

## Download managers

Unprotected file shares work with download managers: no JavaScript, no cookies,
correct `Range` support, and parallel connections. Serving them from a separate
origin means a shared file can never reach a session. A password-protected link
is the exception: it opens a form you have to submit in a browser, so a download
manager cannot use one.

Each non-`HEAD` request for a shared file is counted as a download, so a
download manager that opens many range connections will show a higher count than
one logical download. Requests that never reach a file are not counted: the
password page, a folder listing, and a ZIP bundle download.

## Revoking

Revoking a link takes effect on the next request: the token is looked up per
request, so new downloads and range requests are refused immediately. A transfer
that is already streaming is not cut off mid-response, and on object storage a
presigned URL stays valid until it expires (15 minutes by default).

Every access is recorded, so you can see whether and when a link was used.
