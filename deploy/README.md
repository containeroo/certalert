# Manifests

use `kubectl -k deploy/ -n NAMESPACE` to deploy certalert.
If you possess certificates with unique requirements that certalert cannot accommodate, you can incorporate an initContainer featuring a customized bash script to extract these certificates. For reference, you can explore the pach-deployment, which serves as an illustrative example.

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
