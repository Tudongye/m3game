#!/usr/bin/env bash

protoc -I . pkg.proto --go_out=plugins=grpc:../..

