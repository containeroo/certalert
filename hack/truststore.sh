#!/bin/bash

# Define an associative array
declare -A certs
certs=( ["regular_certificate"]="changeit" ["chain"]="changeit" ["intermediate"]="changeit" ["root"]="changeit" )

mkdir -p ./tests/certs/truststore
pushd ./tests/certs/truststore

# Iterate over the certificates
for name in "${!certs[@]}"; do
  password=${certs[$name]}

  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkeypair -keyalg RSA -alias $name -keystore ${name}.jks -storepass ${password} -validity 365 -dname "CN=$name" -noprompt

  # Export the certificate to a .crt file
  keytool -exportcert -keystore ${name}.jks -alias $name -file ${name}.crt -storepass ${password} -noprompt

  # Import the certificate to a new truststore
  keytool -importcert -file ${name}.crt -alias $name -keystore ${name}_truststore.jks -storepass ${password} -noprompt
done

# create broken truststore file
echo "broken" > broken.jks

# create file with invalid extension
echo "invalid" > cert.invalid

# create file with no extension
echo "no extension" > no_extension