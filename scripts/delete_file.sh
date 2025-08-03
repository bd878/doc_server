#!/usr/bin/bash

# Delete files

token=${1?:"Usage: read_file.sh token id"}
id=${2?:"Usage: read_file.sh token id"}

curl -v -XDELETE "http://138.124.107.242:80/api/docs/$id?token=$token"