ghenvset ?= build/ghenvset

.PHONY: build
build: clean ${ghenvset}

.PHONY: clean
clean:
	$(RM) -R dist
	$(RM) -R build

${ghenvset}:
	go build -buildmode=pie -o $@ .

.PHONY: test
test:
	go clean -testcache
	go test -v -cover -race ./...
	go test -v -cover -buildmode=pie ./pkg/poiana
