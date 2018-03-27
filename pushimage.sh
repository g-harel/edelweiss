#!/bin/sh

PORT=5000

# docker rm -f registry
# docker run -d -p 5000:$PORT --restart always --name registry registry:2
docker pull ubuntu
docker tag ubuntu localhost:$PORT/ubuntu
docker push localhost:$PORT/ubuntu