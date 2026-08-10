#!/usr/bin/env bash

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  MYIP=$(ipconfig | sed -En 's/IPV4//;s/.* (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p' | head -n 1)
else
    MYIP=$(ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p' | head -n 1)
fi

CONTAINER=selectify
IMAGETAG=selectify:$CONTAINER

docker build -t $IMAGETAG .