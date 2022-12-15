#!/usr/bin/env bash
rm -rf binarys
mkdir binarys

GOOS=linux GOARCH=amd64 go build -o binarys/demo_linux_amd64 service.go
GOOS=linux GOARCH=386 go build -o binarys/demo_linux_386 service.go
GOOS=linux GOARCH=arm go build -o binarys/demo_linux_arm service.go
GOOS=linux GOARCH=arm64 go build -o binarys/demo_linux_arm64 service.go
GOOS=windows GOARCH=amd64 go build -o binarys/demo_windows_amd64.exe service.go
GOOS=windows GOARCH=386 go build -o binarys/demo_windows_386.exe service.go
GOOS=windows GOARCH=arm go build -o binarys/demo_windows_arm.exe service.go
GOOS=windows GOARCH=arm64 go build -o binarys/demo_windows_arm64.exe service.go
GOOS=darwin GOARCH=amd64 go build -o binarys/demo_darwin_amd64 service.go
GOOS=darwin GOARCH=arm64 go build -o binarys/demo_darwin_arm64 service.go

cp -r conf ./binarys/conf