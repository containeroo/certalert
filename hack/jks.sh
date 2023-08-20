#!/bin/bash

# Define a function for the keytool command in Docker
keytool() {
    docker run --rm -v "$(pwd)":/certs openjdk:11-jdk keytool "$@"
}

certs=( "regular" "chain" "intermediate" "root" )

mkdir -p ./tests/certs/jks
pushd ./tests/certs/jks

# Iterate the string array using for loop
for name in "${certs[@]}";do
  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkeypair -keyalg RSA -alias ${name} -keystore /certs/${name}.jks -storepass password -validity 365 -dname "CN=${name}, OU=MyOrganization, O=MyCompany, L=MyCity, ST=MyState, C=MyCountry" -storetype JKS -noprompt
done

# create certificate chain
keytool -genkeypair -keyalg RSA -alias root -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=root" -storetype JKS -noprompt
keytool -genkeypair -keyalg RSA -alias intermediate -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=intermediate" -storetype JKS -noprompt
keytool -genkeypair -keyalg RSA -alias leaf -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=leaf" -storetype JKS -noprompt


# Create broken jks file
echo "broken" > broken.jks

# Create file with invalid extension
echo "invalid" > cert.invalid

# Create pkcs12 file with the keystore tools
keytool -genkey -keyalg RSA -alias pkcs12 -keystore /certs/pkcs12.jks -storepass password -validity 365 -dname "CN=pkcs12" -noprompt

popd
