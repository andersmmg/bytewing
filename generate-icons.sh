#!/usr/bin/env bash
# Regenerate every logo and favicon instance from bytewing_fullsize.svg.
#
# The SVG is the source art. It is rasterised once at 1024px into
# bytewing_fullsize.png, and every output below is that master trimmed and
# scaled to the exact dimensions of each file needed!
#
# Requires rsvg-convert (librsvg2-bin) and ImageMagick.
set -euo pipefail

cd "$(dirname "$0")"

SVG=bytewing_fullsize.svg
BASE=bytewing_fullsize.png

command -v rsvg-convert >/dev/null 2>&1 || {
    echo "error: rsvg-convert is required (librsvg2-bin) — none found in PATH" >&2
    exit 1
}
command -v magick >/dev/null 2>&1 || {
    echo "error: ImageMagick is required (no 'magick' in PATH)" >&2
    exit 1
}
[ -f "$SVG" ] || {
    echo "error: $SVG not found in the project root" >&2
    exit 1
}

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# Rasterise at the SVG's native size.
rsvg-convert -w 1024 -h 1024 "$SVG" -o "$tmp/base.png"
magick "$tmp/base.png" -strip "$BASE"
echo "rasterised $SVG -> $BASE (1024x1024)"

# gen WIDTH HEIGHT FILE
gen() {
    magick "$BASE" -trim +repage \
        -resize "${1}x${2}" -gravity center -background none -extent "${1}x${2}" \
        -strip "$3"
    echo "  $3 (${1}x${2})"
}

echo "regenerating icons from $BASE:"

# The sidebar mark is inlined from the vector to be colored by CSS variables.
bbox="$(magick "$BASE" -format '%@' info:)"   # WxH+X+Y
size="${bbox%%+*}"; offs="${bbox#*+}"
vw="${size%%x*}"; vh="${size#*x}"
vx="${offs%%+*}"; vy="${offs#*+}"

for ink in '#34645b' '#a5d3c5'; do
    grep -q "$ink" "$SVG" || \
        echo "warning: $ink not found in $SVG — the app mark may not follow the theme" >&2
done

sed \
    -e 's/ width="1024" height="1024"//' \
    -e "s/<svg /<svg viewBox=\"$vx $vy $vw $vh\" /" \
    -e 's/#34645b/var(--primary)/g' \
    -e 's/#a5d3c5/var(--muted)/g' \
    "$SVG" > web/src/assets/logo.svg
echo "  web/src/assets/logo.svg (viewBox $vx $vy $vw $vh)"

gen 32  32   web/public/favicon.png
gen 32  32   docs/public/favicon-32.png
gen 180 180  docs/public/apple-touch-icon.png
gen 654 604  docs/public/logo.png
gen 173 160  web/src/assets/logo.png

magick "$BASE" -trim +repage -strip -define icon:auto-resize=16,32,48,64,128,256 \
    docs/public/favicon.ico
echo "  docs/public/favicon.ico (16,32,48,64,128,256)"
