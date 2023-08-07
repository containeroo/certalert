#!/bin/bash

# define an associative array
declare -A certs
certs=( ["with_password"]="password" ["without_password"]="" ["intermediate"]="password" ["root"]="password" ["final"]="password")

mkdir -p ./tests/certs/pem
pushd ./tests/certs/pem

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  password=${certs[$name]}
  echo "Generating $name certificate with password: $password"

  # Generate private key with password if it's set
  if [ -z "${password}" ]; then
    openssl genpkey -algorithm RSA -out ${name}.key
  else
    openssl genpkey -algorithm RSA -out ${name}.key -pass pass:${password}
  fi

  # Generate the self-signed certificate
  openssl req -new -x509 -key ${name}.key -out ${name}.crt -days 365 -subj "/CN=$name"

  # Export the certificate and private key to a PEM file
  if [ -z "${password}" ]; then
    cat ${name}.key ${name}.crt > ${name}.pem
  else
    openssl rsa -in ${name}.key -out ${name}.pem -passin pass:${password}
    cat ${name}.pem ${name}.crt > ${name}.pem
  fi
done

# create broken pem file
echo "broken" > broken.pem

# create pem file with a certificate chain
cat final.crt intermediate.crt root.crt > chain.crt
# create pem file with a certificate chain and private key
cat final.key final.crt intermediate.crt root.crt > chain.pem

popd
