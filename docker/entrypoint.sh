#!/bin/sh
set -eu
cd /src
just build
cp -r ./ui/dist/* /srv/.
