SRCS=$(shell find . -type f -name './cmd/demo/*.go') go.mod go.sum

.PHONY: demo generate test run clean

build: demo

demo: $(SRCS)
	go build ./cmd/demo

generate: $(SRCS)
	go generate ./pkg/...

test:
	go test -v ./...

run: demo
	./demo

clean:
	go clean
