#!/bin/sh

docker build -q --platform linux/amd64 -t flyio/percona-server-mysql-operator:ps-orchestrator .
docker push -q flyio/percona-server-mysql-operator:ps-orchestrator
