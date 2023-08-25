#!/bin/bash

# Define an associative array for certificate names
certs=( "cert1" "cert2" "cert3" )

mkdir -p ./tests/certs/p7
pushd ./tests/certs/p7

# Iterate the string array using for loop
for name in "${certs[@]}"; do
  # Generate private key with password
  openssl genpkey -algorithm RSA -out ${name}.key -pass pass:password
  echo "Generated ${name}.key"

  # Generate the self-signed certificate
  openssl req -new -x509 -key ${name}.key -out ${name}.crt -days 365 -subj "/CN=$name"
  echo "Generated ${name}.crt"

  openssl crl2pkcs7 -nocrl -certfile ${name}.crt -out ${name}.p7b
  echo "Generated ${name}.p7b"
done

# create broken p7 file
echo "broken" > broken.p7b

# create file with invalid extension
echo "invalid" > cert.invalid

# create file with a regular certificate
openssl genpkey -algorithm RSA -out regular.key
openssl req -new -x509 -key regular.key -out regular.pem -subj "/CN=regular"
cat regular.pem regular.pem > regular.p7
echo "Created regular.p7"

# Unknown PEM block type
openssl genpkey -algorithm RSA -out unknown_pk.key
openssl req -new -x509 -key unknown_pk.key -out unknown_pem.crt -days 365 -subj "/CN=unknown_pk"
openssl smime -encrypt -aes256 -binary -in message.txt -outform DER -out message.p7 unknown_pem.crt
echo "Created message.p7"

# No Subject
openssl genpkey -algorithm RSA -out no_subject.key -pass pass:password
openssl req -new -x509 -key no_subject.key -out no_subject.crt -days 365 -subj "/"
openssl crl2pkcs7 -nocrl -certfile no_subject.crt -out no_subject.p7b
echo "Created no_subject.p7b"

popd
