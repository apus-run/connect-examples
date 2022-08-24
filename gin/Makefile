GO ?= go

.PHONY: build
build: server gin client

.PHONY: server
server: gin

.PHONY: gin
gin:
	$(GO) build -o bin/$@-server cmd/server/main.go

.PHONY: client
client:
	$(GO) build -o bin/$@ cmd/$@/main.go

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	$(GO) get -u -t ./... && go mod tidy -v

test:
	@$(GO) test -coverprofile coverage.txt -cover ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

clean:
	rm -rf bin
