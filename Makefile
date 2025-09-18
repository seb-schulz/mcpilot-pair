GO ?= go
BUILD_ARGS ?= -v -ldflags '-w -extldflags "-static"' -tags embedded,production
DESTDIR ?= .

.PHONY: build
build: generate
	CGO_ENABLED=0 $(GO) $@ $(BUILD_ARGS) -o $(DESTDIR)/ ./

.PHONY: run
run:
	$(GO) $@  ./


.PHONY: generate
generate:
	go $@ $(BUILD_ARGS) ./...


.PHONY: test vet
test vet: generate
	go $@ -tags embedded,production ./...


.PHONY: clean
clean:
	go $@ ./
	@rm -rf  __debug_bin*
