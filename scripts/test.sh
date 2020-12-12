#!/bin/bash

if [ ! -d cache/testcase/boards ]; then
  echo "to extract cache/testcase/boards"
  tar -zxvf cache/testcase/boards.tar.gz -C cache/testcase
else
  echo "cache/testcase/boards exists"
fi

go build ./...
gotest -v ./... -cover
ipcrm -S 0x00007ffb
ipcrm -M 0x0000fffe
