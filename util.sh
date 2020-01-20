#!/bin/bash

build() {
  GOOS=linux GOARCH=arm GOARM=5 go build -o pi-temp main.go
}

if [[ ${1} == 'build' ]]; then
  build
fi
