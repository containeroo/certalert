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

## `crt-makr`: A Certificate Extraction Helper for `certalert`

`crt-makr` is a specialized utility designed to assist when `certalert` is unable to directly extract certificates. It aims to simplify the extraction process by automating the search and extraction of `.jks` certificates from a given directory and subsequently converting them into `.crt` format.

**Usage**:

```sh
crt-makr [SOURCE_DIRECTORY] [DESTINATION_DIRECTORY]
```

**Arguments**:

- `SOURCE_DIRECTORY`: The directory containing `.jks` certificate files. The script will search up to one subfolder deep in this directory.
- `DESTINATION_DIRECTORY`: The directory where the extracted `.crt` files will be saved.

**How It Works**:

1. The script iterates over all `.jks` files located in the source directory (and its immediate subdirectories).
2. For each `.jks` file:
   - It derives a password name from the filename by:
     1. Removing the `.jks` extension.
     2. Replacing any hyphens (`-`), dots (`.`), and spaces with underscores (`_`).
     3. Appending `_PASSWORD` to the end.
     4. Converting the whole string to uppercase.
   - It then checks the environment variables for this derived password name.
   - Using the password from the environment variables, it first extracts the certificate and private key into a `.p12` format.
   - It then uses `openssl` to extract and save just the certificate in `.crt` format to the specified destination directory.

**Important**:

- Ensure that the required passwords corresponding to each `.jks` file are available in the environment variables. If a password is missing, the script will terminate.
- Always verify the output to ensure the extraction was successful.
