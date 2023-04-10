.PHONY:

build-cli:
	go build -o ./bin/cli app/cli/cmd/main.go

build-server:
	go build -o ./bin/server app/server/cmd/main.go

run-server: build-server
	./bin/server -cfg=config/config.yml

build-wasm:
	GOOS=js GOARCH=wasm go build -o dist/wasm/main.wasm app/wasm/main.go
	cp app/wasm/static/* dist/wasm/