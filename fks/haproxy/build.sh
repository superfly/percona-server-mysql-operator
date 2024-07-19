#!/bin/sh

docker build -q --platform linux/amd64 -t flyio/percona-server-mysql-operator:ps-haproxy .
docker push -q flyio/percona-server-mysql-operator:ps-haproxy
