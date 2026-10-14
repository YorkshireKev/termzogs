#!/bin/sh
set -eu

project_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
cd "$project_dir"

output_dir="$project_dir/dist"
mkdir -p "$output_dir"

build() {
    goos=$1
    goarch=$2
    platform=$3
    extension=$4
    output="$output_dir/termzogs-$platform-$goarch$extension"

    printf 'Building %s/%s -> %s\n' "$goos" "$goarch" "$output"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -o "$output" .
}

build darwin amd64 macos ""
build darwin arm64 macos ""
build linux amd64 linux ""
build windows amd64 windows .exe
