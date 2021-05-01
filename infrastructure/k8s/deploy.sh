#!/bin/bash

helm install nginx-ingress-v1 nginx-stable/nginx-ingress
kubectl apply -f ./namespaces.yml
kubectl apply -f ./groceryspend/hello-world.yml
kubectl apply -f ./groceryspend/ingress.yml