#!/bin/sh

function fail {
    echo
    echo $1
    echo
    echo FAIL
    exit 1
}

server=http://0.0.0.0:8080
here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "testing on $server"
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"
echo "testing put"
NGXFS_CONF=$here/ngxfs.conf ngxfs -put $here/hi.txt /hi.txt
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 200 || fail "Should have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 200 || fail "Should have found /C/hi.txt"
echo "testing get"
NGXFS_CONF=$here/ngxfs.conf ngxfs /hi.txt | grep -q helo || fail "Didn't find /hi.txt"
echo "testing ls"
NGXFS_CONF=$here/ngxfs.conf ngxfs -ls / | grep -q "hi.txt" || fail "Incorrect ls list"
echo "testing del"
NGXFS_CONF=$here/ngxfs.conf ngxfs -del /hi.txt
curl -I --silent $server/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $server/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $server/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"
echo "cleaning up"
curl -XDELETE http://0.0.0.0:8080/B/
curl -XDELETE http://0.0.0.0:8080/C/
echo "win"
