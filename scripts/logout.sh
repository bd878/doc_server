#!/usr/bin/bash

# Terminates auth session

token=${1?:"Usage: read_file.sh token id"}

curl -v -XDELETE "http://138.124.107.242:80/api/auth/$token"
