#!/bin/bash

# 1. copy to orig
cp /tmp/build_ptt.sh /tmp/build_ptt_new.sh

# 2. remove git clone
sed -i 's/git.*//g' /tmp/build_ptt_new.sh

# 3. remove copying-config
sed -i 's/cp .*//g' /tmp/build_ptt_new.sh

# 4. bmake clean
sed -i 's/bmake all/bmake clean all/g' /tmp/build_ptt_new.sh

# 4. add NOKILLWATERBALL
echo "#define NOKILLWATERBALL" >> /home/bbs/pttbbs/pttbbs.conf

# 5. add HOTBOARDCACHE
echo "#define HOTBOARDCACHE 10" >> /home/bbs/pttbbs/pttbbs.conf

# 6. rebuild
bash /tmp/build_ptt_new.sh
