#!/usr/bin/env bash

protoc -I . *.proto --go_out=plugins=grpc:../..

