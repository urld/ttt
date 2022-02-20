#!/bin/sh

set -e 
MAKE="make -e"
export VERSION=$(git describe --exact-match --tags 2>/dev/null || git log -n1 --pretty='%h')

$MAKE install
$MAKE test

GOOS=linux GOARCH=amd64 $MAKE dist
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc $MAKE dist
