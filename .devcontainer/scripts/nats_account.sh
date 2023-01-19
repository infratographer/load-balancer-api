#!/bin/bash
 # script to dump creds for use in our app

 echo "Dumping NATS user creds file"
 nsc generate creds -a LBAAS -n USER > /tmp/user.creds

 echo "Dumping NATS sys creds file"
 nsc generate creds -a SYS -n sys > /tmp/sys.creds
