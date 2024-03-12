#!/bin/bash

set -o errexit
set -o pipefail

# Check dependencies
array=( "helm" "kubectl" )
for i in "${array[@]}"
do
    command -v $i >/dev/null 2>&1 || { 
        echo >&2 "$i is required"; 
        exit 1; 
    }
done

helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace

kubectl -n ingress-nginx patch cm ingress-nginx-controller \
  -p '{"data": {"allow-snippet-annotations":"true"}}'

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

kubectl create namespace alfresco

helm repo add alfresco https://kubernetes-charts.alfresco.com/stable
helm repo update

GLOBAL_KNOWN_URLS=http://localhost
VALUES="values/version_values.yaml,values/resources_values.yaml,values/community_values.yaml"
helm install acs alfresco/alfresco-content-services \
   --values=${VALUES} \
   --set global.search.sharedSecret={{.Secret}} \
   --set global.known_urls=${GLOBAL_KNOWN_URLS} \
   --atomic \
   --timeout 5m0s \
   --namespace=alfresco