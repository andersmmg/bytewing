# Previews

Files can be previewed in place, in a full-screen overlay, with keyboard
navigation between files in the same folder.

## What can be previewed

- **Images** and **PDF**
- **Video**, with seeking, and **audio**
- **Text**, including code and Markdown, shown **as source**, without syntax
  highlighting or rendering

## The server never parses documents

Everything relies on what the browser renders natively (images, video, audio,
and PDF in a sandboxed iframe), while text is shown as raw source. There is **no
document parser in the backend**, so there is no library-parsing vulnerability
to track, which is the category of problem the author of File Browser gave up
trying to fix.

The PDF viewer runs in an iframe with `allow-scripts` and `allow-same-origin`,
because pdf.js will not render pages without scripting. That frame is loaded
from the content origin, not the application origin, so a script inside it has
no access to your session.

## Object storage

When files live in an S3-compatible bucket, downloads are presigned redirects,
so previews work whether or not you run a second origin.

## Local storage

With local storage, in-place preview needs a separate content origin. In
single-origin mode every file is sent as a download instead. See
[One origin, or two?](/guide/installation#one-origin-or-two).
