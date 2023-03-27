#!/bin/bash
echo "GateApp Build..."
cd gateapp/main
go build -gcflags=all="-N -l" -ldflags=-compressdwarf=false .