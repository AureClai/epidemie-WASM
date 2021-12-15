# Epidemy Modelisation with Go and Webassembly

Live Demo : [https://aureclai.github.io/pages/wasm-epidemy/index.html](https://aureclai.github.io/pages/wasm-epidemy/index.html)

## Installation

- Require GO >= 1.15
- Clone repository

### Launch test server

- Go to `cmd/server/`
- Run `go run .`

### Rebuild App

- On PowerShell (Windows) run `build.ps1`
- On Unix system run `GOOS=js GOARCH=wasm go build -o ./assets/epidemie.wasm ./cmd/wasm/`



