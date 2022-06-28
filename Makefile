
EXAMPLES=./examples/**/
BIN=bin/
PROXY=GOPROXY=proxy.golang.org
URI=github.com/ndaba1/gommander
ARTIFACTS=*.prof *.out *.bench *.exe
BENCH=.bench/

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in your PATH, run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

ifeq (, $(shell which benchstat))
$(warning "could not find benchstat in your PATH. Is benchstat installed?")
endif

all: test lint bench

test:
	go test

bench:
	go test --bench=.

lint:
	golangci-lint run

release:
	$(GOPROXY) go list -m $(URI)@$(VERSION)

coverage:
	go test -coverprofile=coverage.out

reports: coverage
	go tool cover -html=coverage.out

benchcmp: 
	benchstat $(BENCH)old.bench $(BENCH)latest.bench

examples: $(EXAMPLES)
	go build -o $(BIN) $(EXAMPLES)

clean: 
	$(RM) -r $(BIN) $(ARTIFACTS) $(BENCH)

