#!/bin/bash
echo "RoleApp Build..."
cd roleapp/main
go build -tags netgo .