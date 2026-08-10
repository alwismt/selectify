#!/usr/bin/env bash

#SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
#cd "$SCRIPT_DIR"

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  MYIP=$(ipconfig | sed -En 's/IPV4//;s/.* (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p' | head -n 1)
else
    MYIP=$(ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p' | head -n 1)
fi


CONTAINER=selectify
IMAGETAG=selectify:$CONTAINER

docker build -t $IMAGETAG .

docker stop $CONTAINER || true && docker rm $CONTAINER || true

#docker run -d --name $CONTAINER --add-host=host.docker.internal:host-gateway -p 80:80 -p 443:443 -p 443:443/udp --net="bridge" $IMAGETAG

docker run -d --name $CONTAINER --add-host="host:$MYIP" -p 80:80 -p 443:443 -p 443:443/udp --net="bridge" -i $IMAGETAG