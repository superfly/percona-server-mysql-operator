#!/bin/sh

IMAGE=flyio/percona-server-mysql-operator:fks-23

docker tag perconalab/percona-server-mysql-operator:js-fks $IMAGE
docker push $IMAGE
