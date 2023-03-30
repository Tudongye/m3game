#!/bin/bash
echo "UidApp Build..."
cd uidapp/main
go build -tags netgo .