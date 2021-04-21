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

# 7. update max user to 200000
sed -i 's/.*#define MAX_USERS.*/#define MAX_USERS (200000)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. update max-active user to 100
sed -i 's/.*#define MAX_ACTIVE.*/#define MAX_ACTIVE (100)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. update max board to 8192
sed -i 's/.*#define\s*MAX_BOARD.*/#define MAX_BOARD (8192)/g' /home/bbs/pttbbs/pttbbs.conf

# 8. shm_offset.c
cp /srv/go-pttbbs/c-pttbbs/shm_offset.c /home/bbs/pttbbs/util
cp /home/bbs/pttbbs/util/Makefile /home/bbs/pttbbs/util/Makefile.c-pttbbs
sed -i 's/initbbs/initbbs shm_offset/g' /home/bbs/pttbbs/util/Makefile

# 9. rebuild
bash /tmp/build_ptt_new.sh
