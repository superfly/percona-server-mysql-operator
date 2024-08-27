#!/bin/sh

IMAGE=flyio/percona-server-mysql-operator:js-fks

docker tag perconalab/percona-server-mysql-operator:js-fks $IMAGE
docker push $IMAGE

run() {
    echo "Building $1"
    local service=$1
    cd $service && ./build.sh &
}

run mysql
run haproxy
run orchestrator

wait
