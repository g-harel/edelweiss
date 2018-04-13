#!/bin/sh

PORT=5000

# docker rm -f registry
# docker run -d -p 5000:$PORT --restart always --name registry registry:2
docker pull alpine
docker tag alpine localhost:$PORT/alpine
docker push localhost:$PORT/alpine
