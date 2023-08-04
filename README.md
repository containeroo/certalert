# certalert

This program serves as a dynamic tool for the purpose of handling and monitoring certificates, and presenting their expiration dates in an epoch format. This design allows Prometheus to poll these metrics.

On invocation of the `/metrics` endpoint by Prometheus, the tool performs a real-time check on the expiration dates of the certificates.

Additionally, `certalert` also supports forwarding the expiration date epoch directly to a Pushgateway server, offering flexibility and control over monitoring workflows.

## Certificate Management

Certificates can be defined with properties such as their `name`, `path`, `type`, and an optional `password`. You have the flexibility to enable or disable specific certificate checks. Additionally, the type of certificate can either be manually defined or determined by the system based on the file extension.

Credentials, such as `passwords`, can be specified in multiple ways: `plain text`, an `environment variable`, or a `file` containing the credentials. For files with multiple key-value pairs, a specific key can be chosen by appending `:{KEY}` at the end of the file path. See `Providing Credentials` for more details.

## Pushgateway Interaction

The program can interact with a Pushgateway server, for which the `address` and `job` label can be defined. It also provides two authentication methods - `Basic` and `Bearer`. For `Basic` authentication, a `username` and `password` are required. For `Bearer` authentication, a `token` is required.

Just like the certificate password, these credentials can also be provided as `plain text`, from an `environment variable`, or from a `file`. See `Providing Credentials` for more details.

## Configuration

### Pushgateway

Below are the available properties for the `Pushgateway` and its nested types:

- **Pushgateway**
  - **Address**: This property specifies the URL of the Pushgateway server.
  - **Job**: This property defines the job label to be attached to pushed metrics.
  - **Auth**: This nested structure holds the authentication details needed for the Pushgateway server. It supports two types of authentication: `Basic` and `Bearer`.

- **Auth**
  - **Basic**: This nested structure holds the basic authentication details.
    - **Username**: This is the username used for basic authentication.
    - **Password**: This is the password used for basic authentication.
  - **Bearer**: This nested structure holds the bearer authentication details.
    - **Token**: This is the bearer token used for bearer authentication.

Please ensure each property is correctly configured to prevent any unexpected behaviors. Remember to provide necessary authentication details under the `Auth` structure based on the type of authentication your Pushgateway server uses.

### Certificate

Here are the available properties for the certificate:

- **Name**: This refers to the unique identifier of the certificate. It's used for distinguishing between different certificates. If not provided, it defaults to the certificate's filename, replacing all spaces (` `), dots (`.`) and underlines (`_`) with a dash (`-`).
- **Enabled**: This toggle enables or disables this check. By default, it is set to `true`.
- **Path**: This specifies the location of the certificate file in your system.
- **Type**: This denotes the type of the certificate. If it's not explicitly specified, the system will attempt to determine the type based on the file extension. Allowed types are: p12, pkcs12, pfx, pem, crt and jks.
- **Password**: This optional property allows you to set the password for the certificate.

#### Providing Credentials

Credentials such as passwords or tokens can be provided in one of the following formats:

- **Plain Text**: Simply input the credentials directly in plain text.
- **Environment Variable**: Use the `env:` prefix, followed by the name of the environment variable that stores the credentials.
- **File**: Use the `file:` prefix, followed by the path of the file that contains the credentials. The file should contain only the credentials.

    In case the file contains multiple key-value pairs, the specific key for the credentials can be selected by appending `:{KEY}` to the end of the path. Each key-value pair in the file must follow the `key = value` format. The system will use the value corresponding to the specified `{KEY}`.

Make sure each credential property is correctly configured to prevent any unexpected behaviors.

__Example__

```yaml
---
pushgateway:
  address: http://pushgateway.monitoring.svc.cluster.local:9091
  job: certalert
certs:
  - name: PEM - without_password
    enabled: true
    path: /certs/pem/without_password_certificate.pem
  - name: PEM - chain
    enabled: true
    path: /certs/pem/chain_certificate.pem
    password: file:/certs/pem/chain_certificate.password
  - name: P12 - with_password
    enabled: true
    path: /certs/p12/with_password_certificate.p12
    password: env:P12_PASSWORD
  - name: P12 - chain
    enabled: true
    path: /certs/p12/chain_certificate.p12
    password: file:/certs/certalert.passwords:{p12_password}
  - name: jks - regular
    enabled: true
    path: /certs/jks/regular.jks
    password: env:JKS_PASSWORD
  - name: jks - chain
    enabled: true
    path: /certs/jks/chain.jks
    password: file:/certs/certalert.passwords:{jks_password}
```

## Available Endpoints

certalert provides the following web-accessible endpoints:

| Endpoint   | Purpose                                                                             |
| :--------- | :---------------------------------------------------------------------------------- |
| `/`        | Fetches and displays all the certificates in a tabular format                       |
| `/config`  | Provides the currently active configuration file. Plaintext passwords are redacted. |
| `/metrics` | Delivers metrics for Prometheus to scrape                                           |
