#!/usr/bin/env bash

protoc -I . -I ../../proto *.proto --go_out=plugins=grpc:../
