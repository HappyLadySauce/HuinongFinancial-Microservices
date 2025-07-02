#!/bin/bash
export KUBECONFIG=/root/.kube/karmada.config

kubectl delete -f appuser.yaml
kubectl delete -f lease.yaml
kubectl delete -f leaseproduct.yaml
kubectl delete -f loan.yaml
kubectl delete -f loanproduct.yaml
kubectl delete -f oauser.yaml
kubectl delete -f app-ingress.yaml
kubectl delete -f oa-ingress.yaml