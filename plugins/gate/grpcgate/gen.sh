#!/usr/bin/env bash

protoc -I . -I ../ -I ../../../meta/metapb *.proto --go_out=../../../.. --go-grpc_out=../../../..
