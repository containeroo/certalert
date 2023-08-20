#!/bin/bash

# Define a function for the keytool command in Docker
keytool() {
    docker run --rm -v "$(pwd)":/certs openjdk:11-jdk keytool "$@"
}

# Define an associative array
declare -A certs
certs=( ["regular"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/truststore
pushd ./tests/certs/truststore

# Iterate over the certificates
for name in "${!certs[@]}"; do

  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkeypair -keyalg RSA -alias $name -keystore /certs/${name}_keystore.jks -storepass password -validity 365 -dname "CN=$name" -noprompt

  # Export the certificate to a .crt file
  keytool -exportcert -keystore /certs/${name}_keystore.jks -alias $name -file /certs/${name}.crt -storepass password -noprompt

  # Create a new truststore and import the certificate
  keytool -importcert -file /certs/${name}.crt -alias $name -keystore /certs/${name}.jks -storepass password -noprompt
done

# Create certificate chain
keytool -importcert -file /certs/root.crt -alias root -keystore /certs/chain.jks -storepass password -noprompt
keytool -importcert -file /certs/intermediate.crt -alias intermediate -keystore /certs/chain.jks -storepass password -noprompt
keytool -importcert -file /certs/regular.crt -alias regular -keystore /certs/chain.jks -storepass password -noprompt

popd

# Create broken truststore file
echo "broken" > ./tests/certs/truststore/broken.jks

# Create file with invalid extension
echo "invalid" > ./tests/certs/truststore/cert.invalid

# Create file with no extension
echo "no extension" > ./tests/certs/truststore/no_extension
