#!/bin/sh

IMAGE=flyio/percona-server-mysql-operator:js-fks

docker tag perconalab/percona-server-mysql-operator:js-fks $IMAGE
docker push $IMAGE

cd mysql && ./build.sh && cd ../haproxy && ./build.sh && cd ../orchestrator && ./build.sh
