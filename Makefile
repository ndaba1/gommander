
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

cover:
	go test -cover

examples: $(EXAMPLES)
	go build -o $(BIN) $(EXAMPLES)

clean: 
	$(RM) -r $(BIN)

