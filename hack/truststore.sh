#!/bin/bash

# Define a function for the keytool command in Docker
keytool() {
    docker run --rm -v "$(pwd)":/certs openjdk:11-jdk keytool "$@"
}

# Define an associative array
declare -A certs
certs=( ["regular"]="password" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/truststore
pushd ./tests/certs/truststore

# Iterate over the certificates
for name in "${!certs[@]}"; do

  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkeypair -keyalg RSA -alias $name -keystore /certs/${name}.jks -storepass password -validity 365 -dname "CN=$name" -noprompt

  # Export the certificate to a .crt file
  keytool -exportcert -keystore /certs/${name}.jks -alias $name -file /certs/${name}.crt -storepass password -noprompt

  # Create a new truststore and import the certificate
  keytool -importcert -file /certs/${name}.crt -alias $name -keystore /certs/${name}_truststore.jks -storepass password -noprompt
done

# Create certificate chain
keytool -genkeypair -keyalg RSA -alias root -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=root" -noprompt
keytool -genkeypair -keyalg RSA -alias intermediate -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=intermediate" -noprompt
keytool -genkeypair -keyalg RSA -alias leaf -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=leaf" -noprompt

# Export the certificates for the chain and create truststore
for name in root intermediate leaf; do
  keytool -exportcert -keystore /certs/chain.jks -alias $name -file /certs/${name}.crt -storepass password -noprompt

  # Import the certificate to the chain truststore
  keytool -importcert -file /certs/${name}.crt -alias $name -keystore /certs/chain_truststore.jks -storepass password -noprompt
done

popd

# Create broken truststore file
echo "broken" > ./tests/certs/truststore/broken.jks

# Create file with invalid extension
echo "invalid" > ./tests/certs/truststore/cert.invalid

# Create file with no extension
echo "no extension" > ./tests/certs/truststore/no_extension
