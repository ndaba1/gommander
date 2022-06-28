
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

fmt:
	$(info *********** running format checks ***********)
	go fmt

test:
	$(info *********** running tests ***********)
	go test

bench:
	$(info *********** running benches ***********)
	go test --bench=.

lint:
	$(info *********** running linting ***********)
	golangci-lint run

release:
	$(GOPROXY) go list -m $(URI)@$(VERSION)

coverage:
	$(info *********** checking coverage ***********)
	go test -coverprofile=coverage.out

reports: coverage
	go tool cover -html=coverage.out

benchcmp: 
	$(info *********** comparing benches ***********)
	benchstat $(BENCH)old.bench $(BENCH)latest.bench

examples: $(EXAMPLES)
	$(info *********** generating example binaries ***********)
	go build -o $(BIN) $(EXAMPLES)

clean: 
	$(info *********** cleaning build and test artifacts ***********)
	$(RM) -r $(BIN) $(ARTIFACTS) $(BENCH)

