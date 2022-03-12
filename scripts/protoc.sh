#!/bin/bash

protoc -I=proto --go_out=. --go_opt=Mrune.proto=./types/proto ./proto/rune.proto
