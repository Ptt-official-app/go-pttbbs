#!/bin/bash
#
# This is the script to rebuild pttbbs in Dockerfile.

# https://github.com/bbsdocker/imageptt/blob/master/pttbbs_conf

# 1. copy to orig
cp /tmp/build_ptt.sh /tmp/build_ptt_new.sh

# 2. remove git clone
sed -i 's/git.*//g' /tmp/build_ptt_new.sh

# 3. remove copying-config
sed -i 's/cp .*//g' /tmp/build_ptt_new.sh

# 4. bmake clean
sed -i 's/pmake all/pmake clean all/g' /tmp/build_ptt_new.sh

# 4.1. boardd
sed -i '/^pmake install/a ## boardd\ncd ${BBSHOME}/pttbbs/daemon/boardd\npmake clean all\ncp ${BBSHOME}/pttbbs/daemon/boardd/boardd ${BBSHOME}/bin' /tmp/build_ptt_new.sh

# 4.2. mand
sed -i '/^pmake install/a ## mand\ncd ${BBSHOME}/pttbbs/daemon/mand\npmake clean all\ncp ${BBSHOME}/pttbbs/daemon/mand/mand ${BBSHOME}/bin' /tmp/build_ptt_new.sh

# 5. add NOKILLWATERBALL
echo "#define NOKILLWATERBALL" >> /home/bbs/pttbbs/pttbbs.conf

# 6. add HOTBOARDCACHE
echo "#define HOTBOARDCACHE (128)" >> /home/bbs/pttbbs/pttbbs.conf

# 7. add USE_EDIT_HISTORY
echo "#define USE_EDIT_HISTORY" >> /home/bbs/pttbbs/pttbbs.conf

# 7. add SAFE_ARTICLE_DELETE
echo "#define SAFE_ARTICLE_DELETE" >> /home/bbs/pttbbs/pttbbs.conf

# 7. add contact-email
echo "#define USEREC_EMAIL_IS_CONTACT" >> /home/bbs/pttbbs/pttbbs.conf
echo "#define ALLOW_REGISTER_WITH_ONLY_CONTACT_EMAIL" >> /home/bbs/pttbbs/pttbbs.conf
echo "#define REQUIRE_CONTACT_EMAIL_TO_CHANGE_PASSWORD" >> /home/bbs/pttbbs/pttbbs.conf

# 7. update max user to 2000000
sed -i 's/.*#define MAX_USERS.*/#define MAX_USERS (2000000)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. update max-active user to 512
sed -i 's/.*#define MAX_ACTIVE.*/#define MAX_ACTIVE (512)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. update max board to 20000
sed -i 's/.*#define\s*MAX_BOARD.*/#define MAX_BOARD (20000)/g' /home/bbs/pttbbs/pttbbs.conf

# 7. add EXTENDED_INCHAR_ANSI
echo "#define EXTENDED_INCHAR_ANSI" >> /home/bbs/pttbbs/pttbbs.conf

# 8. shm_offset.c
cp /srv/go-pttbbs/c-pttbbs/shm_offset.c /home/bbs/pttbbs/util
cp /home/bbs/pttbbs/util/Makefile /home/bbs/pttbbs/util/Makefile.c-pttbbs
sed -i 's/initbbs/initbbs shm_offset/g' /home/bbs/pttbbs/util/Makefile

# 9. rebuild
bash /tmp/build_ptt_new.sh
