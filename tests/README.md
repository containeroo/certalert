# Certificates

To create new certificates for tests you can use one of the following scripts:

- `hack/jks.sh`
- `hack/p7.sh`
- `hack/p12.sh`
- `hack/pem.sh`
- `hack/truststore.sh`

You need `openssl` and `docker` installed (`docker` is for `keytool`).

## ConfigMaps

A `ConfigMap` with all certificates can be created with kustomize:

```sh
kustomize build tests/ -o deploy/certs.yaml
```
