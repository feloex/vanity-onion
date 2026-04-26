#!/bin/bash
echo "Building wasm..."

# find wasm_exec.js, if not found download it from the go repo
if [ ! -f web/wasm_exec.js ]; then
    echo "Downloading wasm_exec.js..."
    curl -s -o web/wasm_exec.js https://raw.githubusercontent.com/golang/go/refs/heads/master/lib/wasm/wasm_exec.js
else
    echo "wasm_exec.js exists"
fi


if [ ! -f web/wasm_exec.js ]; then
    echo "failed to download wasm_exec.js, exiting"
    exit 1
fi


GOOS=js GOARCH=wasm go build -o web/main.wasm .

echo "Server running on http://localhost:8080"

cd web && python3 -m http.server 8080