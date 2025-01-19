#!/bin/bash
# This script provides an example of initializing SYSOP.

# 1. create SYSOP
echo -e "\x1b[1;32m[INFO]\x1b[m 1. To create SYSOP:"

curl -s -X POST -H 'X-Forwarded-For: 127.0.0.1' -d '{"client_info": "test_client_id", "username": "SYSOP", "password": "123123", "password_confirm": "123123", "over18": true}' 'localhost:3456/v1/register'|jq

# 2. login as SYSOP and get access token
echo -e "\x1b[1;32m[INFO]\x1b[m 2. To login as SYSOP and get access token:"

access_token=`curl -s -X POST -H 'X-Forwarded-For: 127.0.0.1' -d '{"client_info": "test_client_id", "username": "SYSOP", "password": "123123"}' 'localhost:3456/v1/token'|jq -c '.access_token'|sed 's/^"//g'|sed 's/"$//g'`

echo -e "\x1b[1;32m[INFO]\x1b[m    access_token: '${access_token}'"

# 3. set permission
echo -e "\x1b[1;32m[INFO]\x1b[m 3. To set permission as PERM_DEFAULT | PERM_POST | PERM_LOGINOK | PERM_MAILLIMIT | PERM_CLOAK | PERM_SEECLOAK | PERM_XEMPT | PERM_SYSOPHIDE | PERM_BM | PERM_ACCOUNTS | PERM_CHATROOM | PERM_BOARD | PERM_SYSOP | PERM_BBSADM:"

curl -s -X POST -H 'X-Forwarded-For: 127.0.0.1' -H "Authorization: bearer ${access_token}" -d '{"perm": 65535}' 'http://localhost:3456/v1/admin/user/SYSOP/setperm'|jq

# 4. show SYSOP info
echo -e "\x1b[1;32m[INFO]\x1b[m 4. To show SYSOP information:"

curl -s -X GET -H 'X-Forwarded-For: 127.0.0.1' -H "Authorization: bearer ${access_token}" 'http://localhost:3456/v1/user/SYSOP/information'|jq
