#!/bin/bash
# Quickly initializing the .PASSWDS using /dev/zero

if [ "$#" != "2" ]; then
    echo "usage: initpasswd.sh [BBSHOME] [N_USER]"
    exit 255
fi

BBSHOME=$1
N_USER=$2
echo "BBSHOME: ${BBSHOME} N_USER: ${N_USER}"

bytes=`expr ${N_USER} \* 512`
dd if=/dev/zero of="${BBSHOME}/.PASSWDS" bs=${bytes} count=1
