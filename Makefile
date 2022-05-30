
EXAMPLES=./examples/**/
BIN=bin/
PROXY=GOPROXY=proxy.golang.org

test:
	go test

bench:
	go test --bench=.

cover:
	go test -cover

examples: $(EXAMPLES)
	go build -o $(BIN) $(EXAMPLES)

clean: 
	$(RM) -r $(BIN)

