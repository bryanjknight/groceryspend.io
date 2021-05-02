#!/bin/bash

helm install nginx-ingress-v1 nginx-stable/nginx-ingress
kubectl apply -f ./namespaces.yml

# TODO: add our services here
# kubectl apply -f ./groceryspend/predict.yml
# kubectl apply -f ./groceryspend/server.yml

# TODO: configure ingress
# kubectl apply -f ./groceryspend/ingress.yml