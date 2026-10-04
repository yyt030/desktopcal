#!/bin/bash
set -e
cd /Users/yueyt/workspaces/desktopcal/build/darwin
mkdir -p icon.iconset
gen() { sips -z "$1" "$1" icon_source.png --out "icon.iconset/$2" > /dev/null; }
gen 16   icon_16x16.png
gen 32   icon_16x16@2x.png
gen 32   icon_32x32.png
gen 64   icon_32x32@2x.png
gen 128  icon_128x128.png
gen 256  icon_128x128@2x.png
gen 256  icon_256x256.png
gen 512  icon_256x256@2x.png
gen 512  icon_512x512.png
gen 1024 icon_512x512@2x.png
iconutil -c icns icon.iconset -o icon.icns
rm -rf icon.iconset
cp icon_source.png ../appicon.png
ls -l icon.icns
