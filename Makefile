
EXAMPLES=./examples/**/
BIN=bin/
PROXY=GOPROXY=proxy.golang.org
URI=github.com/ndaba1/gommander

test:
	go test

bench:
	go test --bench=.

release:
	$(GOPROXY) go list -m $(URI)@$(VERSION)

coverage:
	go test -coverprofile=coverage.out

reports: coverage
	go tool cover -html=coverage.out

profiles: 
	go test -bench=. -run=^# -benchmem -cpuprofile cpu.prof -memprofile mem.prof -benchtime=5s > 0.bench

examples: $(EXAMPLES)
	go build -o $(BIN) $(EXAMPLES)

clean: 
	$(RM) -r $(BIN)

