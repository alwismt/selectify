#!/usr/bin/env bash

## New cert1 for local development
#certbot certonly \
#  --standalone \
#  -d selectify.alwis.dev \
#  --config-dir ./cert \
#  --work-dir ./cert/work \
#  --logs-dir ./cert/logs

## Renew local development cert1
certbot renew \
  --config-dir ./cert \
  --work-dir ./cert/work \
  --logs-dir ./cert/logs