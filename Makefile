run: build
	./cmd/beam/beam
build: deps
	go build -o ./cmd/beam/beam ./cmd/beam
dev:
	go run ./cmd/beam --dir="/Users/tacherasasi/Downloads/"
deps:
	go mod tidy
clean:
	rm -f ./cmd/beam/beam
.PHONY: run build clean