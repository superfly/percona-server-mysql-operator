#!/bin/sh

docker build --platform linux/amd64 -t flyio/percona-server-mysql-operator:ps-orchestrator .
docker push flyio/percona-server-mysql-operator:ps-orchestrator 
