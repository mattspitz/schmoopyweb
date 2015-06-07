#!/bin/bash

set -e
set -x

# args to pass to ssh to access whatever machine we're targeting
SSH_ARGS=$@

pushd src/schmoopy/schmoopy_server
go clean
go build
popd

rm -rf .build
mkdir .build

cp src/schmoopy/schmoopy_server/schmoopy_server .build
cp -r src/static .build

# SCP out the supervisor config, nginx config, static files, and the Go binary

# reload supervisor, restart supervisor, reload nginx
