#!/bin/bash
# This is the script to compile rune protobuf for golang
# rune protobuf is used between golang-server and the frontend (web, app).

protoc -I=proto --go_out=. --go_opt=Mrune.proto=./types/proto ./proto/rune.proto
