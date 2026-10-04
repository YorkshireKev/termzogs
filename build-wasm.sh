# Builds the web version and matching Go WASM runtime for GitHub Pages.
set -eu

goroot=$(go env GOROOT)
if [ -f "$goroot/lib/wasm/wasm_exec.js" ]; then
    wasm_exec="$goroot/lib/wasm/wasm_exec.js"
else
    wasm_exec="$goroot/misc/wasm/wasm_exec.js"
fi
cp "$wasm_exec" docs/wasm_exec.js
GOOS=js GOARCH=wasm go build -o docs/termzogs.wasm
