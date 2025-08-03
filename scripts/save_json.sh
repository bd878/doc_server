#!/usr/bin/bash

# Save json
# with given token

token=${1?:"Usage: save_json.sh token"}

meta=$(echo -n '{\"token\":\"%TOKEN%\",\"name\":\"save_json.curl\",\"file\":false,\"public\":false,\"mime\":\"application\/json\",\"grant\":[]}' | sed -e "s/%TOKEN%/$token/g")
json='{\"test\":\"data\"}'

echo $meta

cmd=`cat <<HERE
sed -e "s/%META%/$meta/g" \
-e "s/%JSON%/$json/g" ./curl/save_json.curl |
curl -v -K -
HERE`
result=`eval "$cmd"`
echo $result
