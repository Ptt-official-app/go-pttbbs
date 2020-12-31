#!/bin/bash

dirs=("cache/testcase" "ptt/testcase")

if [ -f cache/testcase/boards.tar.gz ]; then
  for eachDir in ${dirs[@]}; do
    echo "eachDir: ${eachDir}"
    if [ ! -d ${eachDir}/boards ]; then
      echo "to extract ${eachDir}"
      tar -zxvf cache/testcase/boards.tar.gz -C "${eachDir}"
    else
      echo "${eachDir}/boards exists"
    fi
  done
fi

go build ./... && gotest -v ./... -cover
ipcrm -S 0x00007ffb
ipcrm -M 0x0000fffe
