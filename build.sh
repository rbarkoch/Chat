#!/bin/bash

# MARK: Arguments =============================================================

# help message
USAGE="usage: $0 [path=bin]"
if [ "$#" -gt 1 ]; then echo "$USAGE" && exit 1; fi

SRC="$(dirname "$0")/src" # the source path
BIN="${1:-bin}" # the binaries path
EXEC="chat" # the executable
out="$BIN/$EXEC" # the output path

# MARK: Build =================================================================

# build for a plat-arch and package it
build() { # usage: build <plat> <arch> <id>
  plat="$1"; arch="$2";
  id="$1-$2"
  output="$out"

  # add .exe for windows
  if [ "$plat" = "windows" ]; then
    output="$output.exe"
  fi

  # build the executable
  echo "building for $id..."
  GOOS=$plat GOARCH=$arch go build -C "$SRC" -o "../$output" .

  # package into an archive
  archive="$out-$id.zip"
  zip -j "$archive" "$output" > /dev/null
}

# MARK: Targets ===============================================================

# linux
build linux arm64
build linux amd64

# macos
build darwin arm64
build darwin amd64

# windows
build windows arm64
build windows amd64

# development
echo "building for debug..."
if [ "$(go env GOOS)" = "windows" ]; then
  go build -C "$SRC" -tags=debug -o "../$out.exe" .
  rm "$out"
else # macos | linux
  go build -C "$SRC" -tags=debug -o "../$out" .
  rm "$out.exe"
fi

echo "builds created in $BIN"
