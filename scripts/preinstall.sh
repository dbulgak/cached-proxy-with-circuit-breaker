#!/usr/bin/env bash

set -e

# user creation
useradd -r -m -s /sbin/nologin cachedproxy

# app dir
mkdir /home/cachedproxy/app
chown -R cachedproxy:cachedproxy /home/cachedproxy/app

# logs
mkdir /var/log/za/cachedproxy
chown -R cachedproxy:cachedproxy /var/log/za/cachedproxy