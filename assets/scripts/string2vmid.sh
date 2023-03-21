#!/usr/bin/env bash

function calculateChecksum {
    local hash="$1"
    local sum=0
    for ((i=1; i<=40; i+=2)); do
        byte=${hash:$i:2}
        sum=$((sum + 16#${byte}))
    done
    echo "$((sum % 10000))"
}

if [ "$#" -ne 1 ]; then
    echo "Usage: ./string2vmid.sh <git-sha>"
    exit 1
fi

gitSha="$1"

hash=$(echo -n "$gitSha" | sha1sum | awk '{print $1}')
checksum=$(calculateChecksum "$hash")

printf "2%04d\n" "$checksum"