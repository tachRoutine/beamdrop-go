run: build
	./cmd/beam/beam
build:
	go build -o ./cmd/beam/beam ./cmd/beam
clean:
	rm -f ./cmd/beam/beam
.PHONY: run build clean