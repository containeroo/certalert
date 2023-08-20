#!/bin/bash

# Define an associative array for certificate names
declare -A certs
certs=( ["cert1"]="password1" ["cert2"]="password2" ["cert3"]="password3" )

mkdir -p ./tests/certs/p7
pushd ./tests/certs/p7

# Iterate the string array using for loop
for name in "${!certs[@]}"; do
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

  # Bundle certificates into a PKCS#7 file (no private keys)
  # Here, we're just using the same certificate as an example. In real-world scenarios, you'd bundle different certificates.
  openssl crl2pkcs7 -nocrl -certfile ${name}.crt -out ${name}.p7b
  echo "Generated ${name}.p7b"
done

# create broken p7 file
echo "broken" > broken.p7b

# create file with invalid extension
echo "invalid" > cert.invalid

popd
