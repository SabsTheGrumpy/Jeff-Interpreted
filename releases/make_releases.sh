#!/bin/bash

env GOOS=windows GOARCH=amd64 go build -o ./$1-windows-amd64/
zip -r $1-windows-amd64 $1-windows-amd64

env GOOS=linux GOARCH=amd64 go build -o ./$1-linux-amd64/
tar -czvf $1-linux-amd64.tar.gz $1-linux-amd64

env GOOS=linux GOARCH=arm64 go build -o ./$1-linux-arm64/
tar -czvf $1-linux-arm64.tar.gz $1-linux-arm64


