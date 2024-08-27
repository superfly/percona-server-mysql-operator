#!/bin/sh

set -e

NAME=$1

if [ -z "$NAME" ]; then
  echo "Usage: deploy <cluster-name>"
  exit 1
fi


if kubectl get namespace "$NAME" > /dev/null 2>&1; then
  echo "$NAME namespace exists"
else
  echo "Creating $NAME namespace"
  kubectl create ns $NAME
fi

kubectl apply -f ../deploy/crd.yaml
kubectl apply -f ../deploy/rbac.yaml -n $NAME
kubectl apply -f ../deploy/operator.yaml -n $NAME
kubectl apply -f ../deploy/secrets.yaml -n $NAME
yq eval '.metadata.name = "deploy1"' cluster.yaml | kubectl apply -n deploy1 -f -
