#!/bin/bash
clc -s -e rbset_test.go
cat Version.dat
go mod tidy
go fmt .
echo \* cannot do "staticcheck ." no range func support
# staticcheck .
go vet .
echo \* cannot do "golangci-lint run" no range func support
# golangci-lint run
git st
