#!/bin/bash

SOURCE=$1
DEST=$2

[ -z ${SOURCE} ] && echo "First argument must be the source directory"
[ -z ${DEST} ] && echo "Second argument must be the source directory"

[ ! -d ${SOURCE} ] && echo "Source '${SOURCE}' not found!" && exit 1
[ ! -d ${DEST} ] && echo "Destination '${DEST}' not found!" && exit 1

# remove ending slash
SOURCE=${SOURCE%/}
DEST=${DEST%/}

for jks_file in $(find ${SOURCE} -maxdepth 1 -name "*.jks");do
    base_name=$(basename -s .jks "$jks_file")

    password=$(sed 's/[- \.]/_/g' <<< $base_name)
    export password="${password^^}_PASSWORD"

    [ -z "${!password}" ] && \
      echo "'$password' not found in env vars" && \
      exit 1

    echo "execute: keytool -importkeystore \
            -srckeystore ${SOURCE}/${base_name}.jks \
            -srcstorepass ${!password} \
            -destkeystore /tmp/${base_name}.p12 \
            -deststorepass ${!password} \
            -deststoretype PKCS12 "

    keytool -importkeystore \
            -srckeystore ${SOURCE}/${base_name}.jks \
            -srcstorepass ${!password} \
            -destkeystore /tmp/${base_name}.p12 \
            -deststorepass ${!password} \
            -deststoretype PKCS12

    echo "execute: openssl pkcs12 -in /tmp/${base_name}.p12 \
                -nokeys \
                -out ${DEST}/${base_name}.crt \
                -passin pass:${!password}"

    openssl pkcs12 -in /tmp/${base_name}.p12 \
                -nokeys \
                -out ${DEST}/${base_name}.crt \
                -passin pass:${!password}

    echo "set permissions to ${DEST}/${base_name}.crt (644)"
    chmod 644 ${DEST}/${base_name}.crt

done
