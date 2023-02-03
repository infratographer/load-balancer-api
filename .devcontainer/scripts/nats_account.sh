#!/bin/bash

export PATH="${PATH}:/home/vscode/.nsccli/bin"
export NKEYS_PATH=/nsc/nkeys
export NSC_HOME=/nsc
nsc describe user  -n USER > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "NATS accounts already exist"
    exit 0
fi

# script to dump creds for use in our app
sudo chown -R vscode /nsc

echo "Dumping NATS user creds file"
nsc generate creds -a LBAAS -n USER > /tmp/user.creds

echo "Dumping NATS sys creds file"
nsc generate creds -a SYS -n sys > /tmp/sys.creds
