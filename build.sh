#!/bin/bash

if [ ! -d built ]; then 
    mkdir built
fi

GOOS=linux GOARCH=amd64 go build -o ./built/iredmail-cli-linux