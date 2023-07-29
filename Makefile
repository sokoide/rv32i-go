DEMOSRCS=$(shell find . -type f -name './cmd/demo/*.go') go.mod go.sum
ASMSRCS=$(shell find . -type f -name './cmd/asm/*.go') go.mod go.sum $(YACC_GOS)

# YACC
YACC_DEFS := $(shell find . $(DONT_FIND) -type f -name *.y -print)
YACC_GOS := $(patsubst %.y,%.y.go,$(YACC_DEFS))

.PHONY: demo asm generate test run runasm install-goyacc clean

build: demo asm

demo: $(DEMOSRCS)
	go build ./cmd/demo

asm: $(ASMSRCS) yacc
	go build ./cmd/asm

yacc: $(YACC_GOS)

%.y.go: %.y
	goyacc -p $(basename $(notdir $<)) -o $@ $<
	# if you want to remove '/line...'
	# sed -i.back '/^\/\/line/ d' $@
	# rm ${@}.back

generate: $(DEMOSRCS)
	go generate ./pkg/...

test:
	go test -v ./...

run: demo
	./demo

runasm: asm
	./asm

install-goyacc:
	go install golang.org/x/tools/cmd/goyacc@latest

clean:
	go clean
	rm asm demo
