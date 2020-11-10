#!/bin/sh

sa_jwt=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
ca_crt=$(cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt)

vault auth enable kubernetes &> /dev/null
vault write auth/kubernetes/config \
  token_reviewer_jwt="$sa_jwt" \
  kubernetes_host="https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT" \
  kubernetes_ca_cert="$ca_crt"

vault policy write obscurus /obscurus.hcl

vault write auth/kubernetes/role/obscurus \
  bound_service_account_names=obscurus \
  bound_service_account_namespaces=obscurus \
  policies=obscurus \
  ttl=24h
