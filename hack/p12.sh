#!/bin/bash

# define an associative array
declare -A certs
certs=( ["with_password"]="password" ["without_password"]="" ["chain"]="password" ["intermediate"]="password" ["root"]="password" )

mkdir -p ./tests/certs/p12
pushd ./tests/certs/p12

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  password=${certs[$name]}

  # Generate private key with password if it's set
  if [ -z "${password}" ]; then
    openssl genpkey -algorithm RSA -out ${name}_private_key.key
  else
    openssl genpkey -algorithm RSA -out ${name}_private_key.key -pass pass:${password}
  fi

  # Generate the self-signed certificate
  openssl req -new -x509 -key ${name}_private_key.key -out ${name}_self_signed_certificate.crt -days 365 -subj "/CN=$name"

  # Export the certificate and private key to a PKCS12 (.p12) file
  if [ -z "${password}" ]; then
    openssl pkcs12 -export -out ${name}_certificate.p12 -inkey ${name}_private_key.key -in ${name}_self_signed_certificate.crt -password pass:
  else
    openssl pkcs12 -export -out ${name}_certificate.p12 -inkey ${name}_private_key.key -in ${name}_self_signed_certificate.crt -password pass:${password}
  fi
done

# create broken p12 file
echo "broken" > broken_certificate.p12


# create p12 file with a certificate chain
cat chain_self_signed_certificate.crt intermediate_self_signed_certificate.crt root_self_signed_certificate.crt > chain.crt
openssl pkcs12 -export -out chain_certificate.p12 -inkey chain_private_key.key -in chain_self_signed_certificate.crt -certfile chain.crt -password pass:password

popd
