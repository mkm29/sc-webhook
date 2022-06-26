#!/bin/bash

openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 365 -key ca.key \
  -subj "/C=AU/CN=security-webhook"\
  -out ca.crt

openssl req -newkey rsa:2048 -nodes -keyout server.key \
  -subj "/C=AU/CN=security-webhook" \
  -out server.csr

openssl x509 -req \
  -extfile <(printf "subjectAltName=DNS:security-webhook.default.svc") \
  -days 365 \
  -in server.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt

echo
echo ">> Generating kube secrets..."
kubectl create secret tls security-webhook-tls \
  --cert=server.crt \
  --key=server.key \
  --dry-run=client -o yaml \
  > ./manifests/webhook/webhook.tls.secret.yaml

ca_bundle=$(cat ca.crt | base64 | tr -d '\n')
tls_crt=$(cat server.crt | base64 | tr -d '\n')
tls_key=$(cat server.key | base64 | tr -d '\n')


echo
echo ">> Updating webhook configurations..."
sed -i "s/_CABUNDLE_/$ca_bundle/g" ./manifests/cluster-config/mutating.config.yaml
sed -i "s/_CABUNDLE_/$ca_bundle/g" ./manifests/cluster-config/validating.config.yaml

echo
echo ">> Generating kube secrets..."
sed -i "s/_TLSCRT_/$tls_crt/g" ./manifests/webhook/webhook.tls.secret.yaml
sed -i "s/_TLSKEY_/$tls_key/g" ./manifests/webhook/webhook.tls.secret.yaml

rm ca.crt ca.key ca.srl server.crt server.csr server.key
