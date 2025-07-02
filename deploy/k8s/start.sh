#!/bin/bash
export KUBECONFIG=/root/.kube/karmada.config

kubectl apply -f appuser.yaml
kubectl apply -f lease.yaml
kubectl apply -f leaseproduct.yaml
kubectl apply -f loan.yaml
kubectl apply -f loanproduct.yaml
kubectl apply -f oauser.yaml
kubectl apply -f app-ingress.yaml
kubectl apply -f oa-ingress.yaml