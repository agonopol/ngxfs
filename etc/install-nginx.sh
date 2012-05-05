#!/bin/sh

set -e
set -x

wget http://nginx.org/download/nginx-1.2.0.tar.gz
tar -xvf nginx-1.2.0.tar.gz
rm nginx-1.2.0.tar.gz
cd nginx-1.2.0
# Disable http_rewrite_module to remove pcre dependency
./configure --with-http_dav_module --without-http_rewrite_module
make
make install
cd ../
mkdir -p /tmp/nginx
mv /usr/local/nginx/conf/nginx.conf /usr/local/nginx/conf/nginx.conf.orig
cp ./nginx.conf /usr/local/nginx/conf/
cd /usr/local/nginx/
sbin/nginx
curl -I http://0.0.0.0:8080/