#!/usr/bin/env bash

# clean
rm *.db
go clean

# build
go build main.go

# run
./main

# assert
ls -alth *.db
#cat *.db
