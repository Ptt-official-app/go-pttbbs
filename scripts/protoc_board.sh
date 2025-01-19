#!/bin/bash
# This is the script to compile boardd protobuf for golang

mkdir -p daemon/boardd
protoc -I=pttbbs/daemon/boardd --go_out=. --go_opt=Mboard.proto=./boardd --go-grpc_out=boardd --go-grpc_opt=Mboard.proto=./boardd ./pttbbs/daemon/boardd/board.proto
