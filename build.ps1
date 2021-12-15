$Env:GOOS = "js"
$Env:GOARCH = "wasm"
echo "Building..."
go build -o ./assets/epidemie.wasm ./cmd/wasm/