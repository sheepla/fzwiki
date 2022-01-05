NAME = fzwiki
BIN := bin/$(NAME)
COVERAGE_OUT := .test/cover.out
COVERAGE_HTML := .test/cover.html

.PHONY: build
build:
	go build -o $(BIN)

.PHONY: test
test:
	go test -coverprofile=$(COVERAGE_OUT) ./...

.PHONY: coverage
coverage:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

.PHONY: clean
clean:
	rm $(BIN)
	rm $(COVERAGE_OUT)
	rm $(COVERAGE_HTML)

