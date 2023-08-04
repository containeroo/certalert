#!/bin/bash

# define an associative array
declare -A certs
certs=( ["regular_certificate"]="password" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/jks
pushd ./tests/certs/jks

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  # Generate the key pair and self-signed certificate, and store them in a JKS file
  keytool -genkey -keyalg RSA -alias $name -keystore ${name}.jks -storepass changeit -validity 365 -dname "CN=$name" -noprompt
done

# create certificate chain
keytool -genkeypair -keyalg RSA -alias root -keystore chain.jks -storepass changeit -validity 365 -dname "CN=root" -noprompt
keytool -genkeypair -keyalg RSA -alias intermediate -keystore chain.jks -storepass changeit -validity 365 -dname "CN=intermediate" -noprompt
keytool -genkeypair -keyalg RSA -alias leaf -keystore chain.jks -storepass changeit -validity 365 -dname "CN=leaf" -noprompt

popd
