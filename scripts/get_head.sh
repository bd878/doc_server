#!/usr/bin/bash

# Get file HEAD request

token=${1?:"Usage: get_head.sh token id"}
id=${2?:"Usage: get_head.sh token id"}

curl -v -XHEAD "http://138.124.107.242:80/api/docs/$id?token=$token"
