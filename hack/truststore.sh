#!/bin/bash

# Define an associative array
declare -A certs
certs=( ["regular_certificate"]="password" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/jks
pushd ./tests/certs/jks

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

# For the certificate chain
keytool -genkeypair -keyalg RSA -alias root -keystore chain.jks -storepass changeit -validity 365 -dname "CN=root" -noprompt
keytool -genkeypair -keyalg RSA -alias intermediate -keystore chain.jks -storepass changeit -validity 365 -dname "CN=intermediate" -noprompt
keytool -genkeypair -keyalg RSA -alias leaf -keystore chain.jks -storepass changeit -validity 365 -dname "CN=leaf" -noprompt

# Export the certificates for the chain
for name in root intermediate leaf; do
  keytool -exportcert -keystore chain.jks -alias $name -file ${name}.crt -storepass changeit -noprompt
  
  # Import the certificate to the chain truststore
  keytool -importcert -file ${name}.crt -alias $name -keystore chain_truststore.jks -storepass changeit -noprompt
done

popd
