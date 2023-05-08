#!/bin/bash
echo "WorldApp Build..."
cd worldapp/main
go build -tags netgo .