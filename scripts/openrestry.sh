#!/bin/bash

cp /usr/local/openresty/nginx/conf/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf.c-pttbbs
sed -i 's/bbsdocker\.github\.io/term.devptt.site/g' /usr/local/openresty/nginx/conf/nginx.conf
