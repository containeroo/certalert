#!/bin/bash

# Define a function for the keytool command in Docker
keytool() {
    docker run --rm -v "$(pwd)":/certs openjdk:11-jdk keytool "$@"
}

# define an associative array
declare -A certs
certs=( ["regular"]="password" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/jks
pushd ./tests/certs/jks

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  password=${certs[$name]}
  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkey -keyalg RSA -alias $name -keystore /certs/${name}.jks -storepass "${password}" -validity 365 -dname "CN=$name" -noprompt
done

# create certificate chain
keytool -genkeypair -keyalg RSA -alias root -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=root" -noprompt
keytool -genkeypair -keyalg RSA -alias intermediate -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=intermediate" -noprompt
keytool -genkeypair -keyalg RSA -alias leaf -keystore /certs/chain.jks -storepass password -validity 365 -dname "CN=leaf" -noprompt

popd

# Create broken jks file
echo "broken" > ./tests/certs/jks/broken.jks

# Create file with invalid extension
echo "invalid" > ./tests/certs/jks/cert.invalid

# Create file with no extension
echo "no extension" > ./tests/certs/jks/no_extension
