#!/usr/bin/bash

# Save file
# with given token

token=${1?:"Usage: save_file.sh token"}

meta=$(echo -n '{\"token\":\"%TOKEN%\",\"name\":\"save_file.curl\",\"file\":true,\"public\":false,\"mime\":\"text\/plain\",\"grant\":[\"test8\",\"test3\"]}' | sed -e "s/%TOKEN%/$token/g")
json='{\"test\":\"data\"}'

echo $meta

cmd=`cat <<HERE
sed -e "s/%META%/$meta/g" \
-e "s/%JSON%/$json/g" ./curl/save_file.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
