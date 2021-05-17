#!/bin/bash

kind_cluster=$(kind get clusters 2>&1)

if [[ $kind_cluster =~ "No kind clusters found" ]]; then
  echo "No kind cluster found..."
  echo "Creating kind cluster..."
  kind create cluster

  echo "Setting up Certificate manager..."
  kubectl create -f config/external/cert-manager.yaml

  while :
  do
    cert_manager_pods=$(kubectl wait --for=condition=Ready --timeout=60s pods --all -n cert-manager 2>&1)
    if [[ $cert_manager_pods =~ "error" ]]; then
      echo "$cert_manager_pods"
      continue
    fi
    echo "$cert_manager_pods"
    break
  done

  echo "Creating Certificates..."

  docker pull busybox
  kind load docker-image busybox
  kubectl create ns system
  kustomize build config/certmanager | kubectl apply -f -
  kubectl create -f hack/busybox-cert.yaml

  while :
  do
    busybox_cert_pod=$(kubectl wait --for=condition=Ready --timeout=60s pod -l app=busybox-cert -n system 2>&1)
    if [[ $busybox_cert_pod =~ "condition met" ]]; then
      echo "$busybox_cert_pod"
      break
    fi
    echo "$busybox_cert_pod"
  done

  echo "Downloading Certificates to local machine..."

  busybox_cert_pod_name=$(kubectl get pod -l app=busybox-cert -n system -o jsonpath="{.items[0].metadata.name}")

  echo "$busybox_cert_pod_name"

  kubectl exec -n "system" "$busybox_cert_pod_name" -- tar cf - "/tmp/k8s-webhook-server/serving-certs" | tar xf - --strip-components 2
  rm -rf /tmp/k8s-webhook-server
  mkdir /tmp/k8s-webhook-server
  mv serving-certs /tmp/k8s-webhook-server

  kubectl delete -f hack/busybox-cert.yaml
  kubectl delete ns system
fi

make install run

#  go clean -testcache && go test ./api/... -v -cover


