#!/bin/bash

ldflags="-s -w"

go env -w GOPROXY=https://goproxy.cn,direct

# Plugin
filepath="plugin"
for item in "$filepath"/*.go; do
  buf=${item%.go}
  name=${buf##*/}
  echo $name
  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$ldflags" -o plugin/filter-"$name" plugin/"$name".go
  upx plugin/filter-"$name"
done

# Test
target="plugin-filter-test"
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$ldflags" -o $target main.go
upx $target
