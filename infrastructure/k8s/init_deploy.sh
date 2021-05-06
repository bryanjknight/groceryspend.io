#!/bin/bash
set -eo pipefail

#### NOTE ####
# Requires kubeconfig setup in ~/.kube/config
##############

# get current directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Namespaces are handled by terraform as it needs the namespaces for the secrets

# run all deploys
find $SCRIPT_DIR -name "*-deploy.yml" -print0 | xargs -n1 -0 kubectl apply -f

# run all svc
find $SCRIPT_DIR -name "*-svc.yml" -print0 | xargs -n1 -0 kubectl apply -f

# run ingress
kubectl apply -f ingress.yml

# TODO: get the external IP of the load balancer