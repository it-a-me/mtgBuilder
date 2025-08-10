default:
  just --list

fetch-cards:
  test ! -f cards.json && go run ./cmd/fetch cards.json || true
  cp cards.json ./ui/public/cards.json

download-deps-go:
  go mod download

[working-directory: 'ui']
download-deps-npm:
  deno install

download-deps: download-deps-go download-deps-npm

[working-directory: 'wasmlib']
build-wasm:
  GOOS=js GOARCH=wasm go build
  cp wasmlib ../ui/public/wasmlib.wasm

[working-directory: 'ui']
build: download-deps fetch-cards build-wasm
  deno run build

[working-directory: 'ui']
build-gh: download-deps fetch-cards build-wasm
  deno run build-gh

[working-directory: 'ui']
run-dev: download-deps fetch-cards build-wasm
  deno run dev || true


