#!/bin/bash

# define an associative array
declare -A certs
certs=( ["with_password"]="password" ["without_password"]="" ["intermediate"]="password" ["root"]="password" ["final"]="password" )

mkdir -p ./tests/certs/p12
pushd ./tests/certs/p12

# Iterate the string array using for loop
for name in "${!certs[@]}";do
  password=${certs[$name]}

  # Generate private key with password if it's set
  if [ -z "${password}" ]; then
    openssl genpkey -algorithm RSA -out ${name}.key
  else
    openssl genpkey -algorithm RSA -out ${name}.key -pass pass:${password}
  fi
  echo "Generated ${name}.key"

  # Generate the self-signed certificate
  openssl req -new -x509 -key ${name}.key -out ${name}.crt -days 365 -subj "/CN=$name"
  echo "Generated ${name}.crt"

  # Export the certificate and private key to a PKCS12 (.p12) file
  if [ -z "${password}" ]; then
    openssl pkcs12 -export -out ${name}.p12 -inkey ${name}.key -in ${name}.crt -password pass:
  else
    openssl pkcs12 -export -out ${name}.p12 -inkey ${name}.key -in ${name}.crt -password pass:${password}
  fi
  echo "Generated ${name}.p12"
done

# create broken p12 file
echo "broken" > broken.p12

# create file with invalid extension
echo "invalid" > cert.invalid

# create file with no extension
echo no_extension > no_extension

# create p12 file with a certificate chain
cat final.crt intermediate.crt root.crt > chain.crt
openssl pkcs12 -export -out chain.p12 -inkey final.key -in final.crt -certfile chain.crt -password pass:password

popd
