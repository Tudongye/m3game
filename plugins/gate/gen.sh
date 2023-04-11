#!/usr/bin/env bash

protoc -I . -I ../../meta/metapb *.proto --go_out=plugins=grpc:../../..
