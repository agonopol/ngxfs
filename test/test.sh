#!/bin/sh

function fail {
    echo $1
    exit 1
}

server=http://0.0.0.0:8080

echo "testing on $server"
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"
echo "testing put"
NGXFS_CONF=./ngxfs.conf ngxfs -put hi.txt /hi.txt
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 200 || fail "Should have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 200 || fail "Should have found /C/hi.txt"
echo "testing get"
NGXFS_CONF=./ngxfs.conf ngxfs /hi.txt | grep -q helo || fail "Didn't find /hi.txt"
echo "testing del"
NGXFS_CONF=./ngxfs.conf ngxfs -del /hi.txt
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"
echo "win"
