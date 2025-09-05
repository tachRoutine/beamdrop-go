run: build
	./cmd/beam/beam
build:
	go build -o ./cmd/beam/beam ./cmd/beam

dev:
	go run ./cmd/beam --dir=./
clean:
	rm -f ./cmd/beam/beam
.PHONY: run build clean