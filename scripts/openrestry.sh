#!/bin/bash
# Reminding how we can update bbs_origin_checked in the nginx.conf for wsproxy.lua.
#
# https://github.com/ptt/pttbbs/blob/master/daemon/wsproxy/wsproxy.lua
# https://github.com/bbsdocker/imageptt/blob/master/nginx_conf_ws#L36

cp /usr/local/openresty/nginx/conf/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf.c-pttbbs
sed -i 's/bbsdocker\.github\.io/term.devptt.dev/g' /usr/local/openresty/nginx/conf/nginx.conf
