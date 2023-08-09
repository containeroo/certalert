# certalert

This program serves as a dynamic tool for the purpose of handling and monitoring certificates, and presenting their expiration dates in an epoch format. This design allows Prometheus to poll these metrics.

On invocation of the `/metrics` endpoint by Prometheus, the tool performs a real-time check on the expiration dates of the certificates.

Additionally, `certalert` also supports forwarding the expiration date epoch directly to a Pushgateway server, offering flexibility and control over monitoring workflows.

## Usage

The primary function is to utilize the `serve` command to initiate a web server that displays metrics for Prometheus to retrieve.

## Global Flags

- `-c, --config`: Specify the path to a config file (Default is `$HOME/.certalert.yaml`).
- `-v, --verbose`: Enable verbose output.
- `-s, --silent`: Enable silent mode, only showing errors.
- `-f, --fail-on-error`: Fail on error.
- `-V, --version`: Print the current version and exit.

### Basic Commands

1. **serve**: Launches a web server to expose certificate metrics.

    ```bash
    certalert serve [flags]
    ```

    Flags:
    - `-H, --hostname`: Specify the hostname to listen on (Default: `localhost`).
    - `-p, --port`: Specify the port to listen on (Default: `8080`).

    Examples:

    ```bash
    # Launch the web server on localhost:8080.
    certalert serve --hostname localhost --port 8080
    ```


2. **print**: Prints information about the certificates.

    ```bash
    certalert print [CERTIFICATE_NAME...] [flags]
    ```

    Flags:
    - `-A, --all`: Prints all certificates.
    - `-o, --output`: Specify the output format. Supported formats: `text`, `json`, `yaml`.

    Examples:

    ```bash
    # Print all certificates in JSON format.
    certalert print --all --output json

    # Print a specific certificate named 'example-cert' in the default format.
    certalert print example-cert
    ```

3. **push**: Push certificate expiration as an epoch to a Prometheus Pushgateway instance.

    ```bash
    certalert push [CERTIFICATE_NAME...] [flags]
    ```

    Flags:
    - `-A, --all`: Push all certificates.
    - `-i, --insecure-skip-verify`: Skip TLS certificate verification.

    Examples:

    ```bash
    # Push metadata for all certificates.
    certalert push --all

    # Push metadata for a specific certificate named 'example-cert'.
    certalert push example-cert
    ```

## Certificate Management

Certificates can be defined with properties such as their `name`, `path`, `type`, and an optional `password`. You have the flexibility to enable or disable specific certificate checks. Additionally, the `type` of certificate can either be manually defined or determined by the system based on the file extension.

Credentials, such as `passwords`, can be specified in multiple ways: `plain text`, an `environment variable`, or a `file` containing the credentials. For files with multiple key-value pairs, a specific key can be chosen by appending `:{KEY}` at the end of the file path. See `Providing Credentials` for more details.

## Pushgateway Interaction

The program can interact with a Pushgateway server, for which the `address` and `job` label can be defined. It also provides two authentication methods - `Basic` and `Bearer`. For `Basic` authentication, a `username` and `password` are required. For `Bearer` authentication, a `token` is required.

Just like the certificate password, these credentials can also be provided as `plain text`, from an `environment variable`, or from a `file`. See `Providing Credentials` for more details.

## Configuration

The certificates must be configured in a file. The config file can be `yaml`, `json` or `toml`. The config file should be loaded automatically if changed. Please check the log output to control if the automatic config reload works in your environment. The endpont `/-/reload` also reloads the configuration.

### Pushgateway

Below are the available properties for the `Pushgateway` and its nested types:

- **pushgateway**
  - **address**: The URL of the Pushgateway server.
  - **job**: The job label to be attached to pushed metrics.
  - **insecureSkipVerify** Skip TLS certificate verification. Defaults to `false`.
  - **auth**: This nested structure holds the authentication details needed for the Pushgateway server. It supports two types of authentication: `Basic` and `Bearer`.

- **auth**
  - **basic**: This nested structure holds the basic authentication details.
    - **username**: Username used for basic authentication.
    - **password**: Password used for basic authentication.
  - **bearer**: This nested structure holds the bearer authentication details.
    - **token**: Bearer token used for bearer authentication.

Please ensure each property is correctly configured to prevent any unexpected behaviors. Remember to provide necessary authentication details under the `Auth` structure based on the type of authentication your Pushgateway server uses.

### Certificate

Here are the available properties for the certificate:

- **name**: This refers to the unique identifier of the certificate. It's used for distinguishing between different certificates. If not provided, it defaults to the certificate's filename, replacing all spaces (` `), dots (`.`) and underlines (`_`) with a dash (`-`).
- **enabled**: This toggle enables or disables this check. By default, it is set to `true`.
- **path**: This specifies the location of the certificate file in your system.
- **type**: This denotes the type of the certificate. If it's not explicitly specified, the system will attempt to determine the type based on the file extension. Allowed types are: p12, pkcs12, pfx, pem, crt and jks.
- **password**: This optional property allows you to set the password for the certificate.

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
  insecureSkipVerify: false
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

| Endpoint    | Purpose                                                                             |
| :---------- | :---------------------------------------------------------------------------------- |
| `/`         | Fetches and displays all the certificates in a tabular format                       |
| `/-/reload` | Reloads the configuration                                                           |
| `/config`   | Provides the currently active configuration file. Plaintext passwords are redacted. |
| `/metrics`  | Delivers metrics for Prometheus to scrape                                           |
| `/healthz`  | Get health status                                                                   |
