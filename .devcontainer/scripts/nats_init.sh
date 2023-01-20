#!/bin/bash
 # script to bootstrap a nats operator environment

 echo "Creating NATS operator"
 nsc add operator --generate-signing-key --sys --name LOCAL
 nsc edit operator -u 'nats://nats:4222'
 nsc list operators
 nsc describe operator

 export OPERATOR_SIGNING_KEY_ID=`nsc describe operator -J | jq -r '.nats.signing_keys | first'`

 echo "Creating NATS account for load-balancer-api"
 nsc add account -n LBAAS -K ${OPERATOR_SIGNING_KEY_ID}
 nsc edit account LBAAS --sk generate --js-mem-storage -1 --js-disk-storage -1 --js-streams -1 --js-consumer -1
 nsc describe account LBAAS

 export ACCOUNTS_SIGNING_KEY_ID=`nsc describe account LBAAS -J | jq -r '.nats.signing_keys | first'`

 echo "Creating NATS user for load-balancer-api"
 nsc add user -n USER -K ${ACCOUNTS_SIGNING_KEY_ID}
 nsc describe user USER

 echo "Generating NATS resolver.conf"
 nsc generate config --mem-resolver --sys-account SYS --config-file /nats/resolver.conf --force
