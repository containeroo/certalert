#!/bin/bash

# define an associative array
declare -A certs
certs=( ["with_password"]="password" ["without_password"]="" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/pem
pushd ./tests/certs/pem

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  password=${certs[$name]}

  # Generate private key with password if it's set
  if [ -z "$password" ]; then
    openssl genpkey -algorithm RSA -out ${name}_private_key.key
  else
    openssl genpkey -algorithm RSA -out ${name}_private_key.key -pass pass:$password
  fi

  # Generate the self-signed certificate
  openssl req -new -x509 -key ${name}_private_key.key -out ${name}_self_signed_certificate.crt -days 365 -subj "/CN=$name"

  # Export the certificate and private key to a PEM file
  if [ -z "$password" ]; then
    cat ${name}_private_key.key ${name}_self_signed_certificate.crt > ${name}_certificate.pem
  else
    openssl rsa -in ${name}_private_key.key -out ${name}_private_key.pem -passin pass:$password
    cat ${name}_private_key.pem ${name}_self_signed_certificate.crt > ${name}_certificate.pem
  fi
done

# create broken pem file
echo "broken" > broken_certificate.pem

# create pem file with a certificate chain
cat chain_self_signed_certificate.crt intermediate_self_signed_certificate.crt root_self_signed_certificate.crt > chain.crt
cat chain_private_key.key chain.crt > chain_certificate.pem

popd
