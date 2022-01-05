NAME = fzwiki

COVER_OUT := .test/cover.out
COVER_HTML := .test/cover.html

.PHONY: build
build:
	go build -o bin/$(NAME)

.PHONY: test
test:
	go test -coverprofile=$(COVER_OUT) ./...

.PHONY: coverage
coverage:
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

