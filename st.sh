#!/bin/bash
clc -s
cat Version.dat
go mod tidy
go fmt .
echo \* cannot do "staticcheck ." no range func support
# staticcheck .
go vet .
echo \* cannot do "golangci-lint run" no range func support
# golangci-lint run
git st
