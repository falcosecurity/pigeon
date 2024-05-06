pigeon ?= pigeon

.PHONY: build
build: clean ${pigeon}

.PHONY: clean
clean:
	$(RM) -R dist
	$(RM) -R ${pigeon}

${pigeon}:
	go build -buildmode=pie -o $@ .

.PHONY: test
test:
	go clean -testcache
	go test -v -cover -race ./...
	go test -v -cover -buildmode=pie ./pkg/pigeon
