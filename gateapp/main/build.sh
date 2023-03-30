#!/bin/bash
echo "GateApp Build..."
cd gateapp/main
go build -tags netgo .