#!/bin/sh
set -eu
cd /src
just build-gh
cp -r ./ui/dist/* /srv/.
