#!/bin/bash

dirs=("cache/testcase" "ptt/testcase")

for eachDir in ${dirs[@]}; do
  echo "eachDir: ${eachDir}"
  if [ ! -d ${eachDir}/boards ]; then
    echo "to extract ${eachDir}"
    tar -zxvf cache/testcase/boards.tar.gz -C "${eachDir}"
  else
    echo "${eachDir}/boards exists"
  fi
done

go build ./...
gotest -v ./... -cover
ipcrm -S 0x00007ffb
ipcrm -M 0x0000fffe
