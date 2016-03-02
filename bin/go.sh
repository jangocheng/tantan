#!/bin/sh

x=`date +%F_%T`
y=`git rev-parse HEAD`

OUT=tantan_linux_amd64
GOOS=linux go build -ldflags "-X main.date=$x -X main.rev=$y" -o $OUT
