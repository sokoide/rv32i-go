TARGET=rv32i-go
SRCS=$(shell find . -type f -name '*.go') go.mod go.sum

.PHONY: $(TARGET) generate test run clean

build: $(TARGET)

$(TARGET): $(SRCS)
	go build

generate: $(SRCS)
	go generate ./pkg/...

test:
	go test -v ./...

run: $(TARGET)
	./$(TARGET)

clean:
	go clean
