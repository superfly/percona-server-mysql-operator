#!/bin/sh

IMAGE=flyio/percona-server-mysql-operator:js-fks

docker tag perconalab/percona-server-mysql-operator:js-fks $IMAGE
docker push $IMAGE
