#!/bin/sh

IMAGE=flyio/percona-server-mysql-operator:fks-24

docker tag perconalab/percona-server-mysql-operator:js-fks $IMAGE
docker push $IMAGE
