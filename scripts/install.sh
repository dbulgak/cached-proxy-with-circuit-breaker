#!/usr/bin/env bash

set -e

# install dependencies
yum -y install redis supervisor

systemctl enable supervisord && systemctl start supervisord
systemctl enable redis && systemctl start redis

# install supervisord configuration
cp ../init/cachedproxy.ini /etc/supervisord.d/
supervisorctl update za_cachedproxy
