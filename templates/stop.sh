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

helm uninstall -n alfresco acs

kubectl delete namespace alfresco

{{- if eq .Kubernetes "docker-desktop" }}
helm uninstall -n ingress-nginx ingress-nginx
{{- end}}

kubectl delete namespace ingress-nginx