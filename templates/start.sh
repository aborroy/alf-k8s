#!/bin/bash

set -o errexit
set -o pipefail

# Check dependencies
array=( "helm" "kubectl" {{ if eq .Kubernetes "kind" }}"kind" {{ end }})
for i in "${array[@]}"
do
    command -v $i >/dev/null 2>&1 || { 
        echo >&2 "$i is required"; 
        exit 1; 
    }
done

{{- if eq .Kubernetes "kind" }}

kind delete cluster

cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

{{- end }}

{{- if eq .Kubernetes "docker-desktop" }}
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
{{- else }}
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
{{- end }}

kubectl -n ingress-nginx patch cm ingress-nginx-controller \
  -p '{"data": {"allow-snippet-annotations":"true"}}'

kubectl -n ingress-nginx rollout status deployment ingress-nginx-controller --timeout=90s

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