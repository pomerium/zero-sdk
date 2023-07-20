//go:build tools
// +build tools

package main

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen" // for oapi-codegen
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
