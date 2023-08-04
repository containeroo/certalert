# CertAlert

CertAlert can handle a variety of certificate types, including .p12, .pkcs12, .pem, .crt, and .jks files.

You can execute specific commands for different actions:

- Use the `push` command to manually push metrics to the Prometheus Pushgateway.
- Use the `serve` command to start a server that provides a `/metrics` endpoint for Prometheus to scrape.
