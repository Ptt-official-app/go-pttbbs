#!/bin/bash

go build ./... && gotest ./... -coverprofile cover.out
ipcrm -S 0x00007ffb
ipcrm -M 0x0000fffe


go tool cover -html=cover.out
