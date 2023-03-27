#!/usr/bin/env bash

protoc -I . -I ../../../meta *.proto --go_out=../ --go-grpc_out=../../ 
