#!/usr/bin/env bash

protoc  -I $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.8/third_party/googleapis -I . -I ../../proto *.proto --go_out=../ --go-grpc_out=../ --grpc-gateway_out=../ 
