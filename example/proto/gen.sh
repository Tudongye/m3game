#!/usr/bin/env bash

protoc -I . -I ../../meta/metapb *.proto --go_out=../ --go-grpc_out=../ --grpc-gateway_out=../ 
