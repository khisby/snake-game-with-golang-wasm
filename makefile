
.PHONY: all run compile

compile:
	GOOS=js GOARCH=wasm go build -o static/main.wasm wasm/main.go

run:
	go run server.go

all: compile run