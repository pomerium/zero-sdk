GOBIN=$(shell pwd)/bin

.PHONY: all
all: test
	@echo "@==> $@"


.PHONY: generate
generate:
	@echo "@==> $@"
	go generate ./...

.PHONY: lint
lint:
	@echo "@==> $@"
	golangci-lint run

.PHONY: proto
proto: proto-tools
	@echo "@==> $@"
	@PATH=$(GOBIN):$(PATH) buf generate

.PHONY: proto-tools
proto-tools:
	@echo "@==> $@"
	@mkdir -p $(GOBIN)
	@GOBIN=$(GOBIN) go install \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/bufbuild/buf/cmd/buf

.PHONY: test
test: 
	@echo "@==> $@"
	go test ./...
