build:
	mkdir dist
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -o dist/chip8 src/main.go
run: build
	./dist/chip8
