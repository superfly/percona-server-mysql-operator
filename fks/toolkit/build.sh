#!/bin/sh

docker build --platform linux/amd64 -t flyio/percona-server-mysql-operator:toolkit-$1 .
docker push flyio/percona-server-mysql-operator:toolkit-$1
