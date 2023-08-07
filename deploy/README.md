# Manifests

## Certificates

To create new certificates you one of the following scripts:

- `hack/pem.sh`
- `hack/p12.sh`
- `hack/jks.sh`

You need `openssl` and `keytool` installed.

## ConfigMaps

`ConfigMaps` with certificates can be created with kustomize:

```sh
cd tests/certs
kustomize build tests/certs -o deploy/certs/
```
