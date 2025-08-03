#!/usr/bin/bash

# List files

token=${1?:"Usage: list.sh token limit"}
limit=${2?:"Usage: list.sh token limit"}

curl -XGET "http://138.124.107.242:80/api/docs?token=$token&key=mime&value=text/plain&limit=$limit&login=test1"
