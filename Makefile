.PHONY: build run clean

build:
	GOOS=js GOARCH=wasm go build -o web/snake.wasm cmd/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js web/

run: build
	go run server/main.go

clean:
	rm -f web/snake.wasm
	rm -f web/wasm_exec.js