#!/bin/sh
set -eu

project_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
cd "$project_dir"

output_dir="$project_dir/dist"
mkdir -p "$output_dir"

version=$(sed -n 's/^const Number = "\(.*\)"$/\1/p' version/version.go)
if [ -z "$version" ]; then
    printf 'Could not read the release version from version/version.go\n' >&2
    exit 1
fi

build() {
    goos=$1
    goarch=$2
    platform=$3
    extension=$4
    output="$output_dir/termzogs-v$version-$platform-$goarch$extension"

    printf 'Building %s/%s -> %s\n' "$goos" "$goarch" "$output"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -o "$output" .
}

build darwin amd64 macos ""
build darwin arm64 macos ""
build linux amd64 linux ""
build windows amd64 windows .exe

goroot=$(go env GOROOT)
if [ -f "$goroot/lib/wasm/wasm_exec.js" ]; then
    wasm_exec="$goroot/lib/wasm/wasm_exec.js"
else
    wasm_exec="$goroot/misc/wasm/wasm_exec.js"
fi

printf 'Building js/wasm -> %s/docs/termzogs.wasm\n' "$project_dir"
cp "$wasm_exec" "$project_dir/docs/wasm_exec.js"
GOOS=js GOARCH=wasm go build -trimpath -o "$project_dir/docs/termzogs.wasm" .
