#!/bin/bash

# 1. copy to orig
cp /tmp/build_ptt.sh /tmp/build_ptt_new.sh

# 2. remove git clone
sed -i 's/git.*//g' /tmp/build_ptt_new.sh

# 3. remove copying-config
sed -i 's/cp .*//g' /tmp/build_ptt_new.sh

# 4. bmake clean
sed -i 's/bmake all/bmake clean all/g' /tmp/build_ptt_new.sh

# 5. add NOKILLWATERBALL
echo "#define NOKILLWATERBALL" >> /home/bbs/pttbbs/pttbbs.conf

# 6. add HOTBOARDCACHE
echo "#define HOTBOARDCACHE (10)" >> /home/bbs/pttbbs/pttbbs.conf

# 7. update max-active user to 100
sed -i 's/.*#define MAX_ACTIVE.*/#define MAX_ACTIVE (100)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. rebuild
bash /tmp/build_ptt_new.sh
