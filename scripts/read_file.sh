#!/usr/bin/bash

# Read file

token=${1?:"Usage: read_file.sh token id"}
id=${2?:"Usage: read_file.sh token id"}

curl -v -XGET "http://138.124.107.242:80/api/docs/$id?token=$token"
