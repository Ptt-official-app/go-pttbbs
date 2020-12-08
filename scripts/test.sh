#!/bin/bash

go build ./...
gotest -v ./... -cover
ipcrm -S 0x00007ffb
ipcrm -M 0x0000fffe
