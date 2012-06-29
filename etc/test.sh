#!/bin/sh

function fail {
    echo
    echo $1
    echo
    echo FAIL
    exit 1
}

host=http://0.0.0.0:8080
here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export NGXFS_CONF=http://0.0.0.0:8080/ngxfs.conf

echo "testing on $host"

echo "ensuring clean host"
curl -I --silent $host/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $host/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $host/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"

echo "testing put"
ngxfs -put $here/hi.txt /hi.txt
curl -I --silent $host/A/hi.txt | grep -q 200 || fail "Should have found /A/hi.txt"
curl -I --silent $host/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $host/C/hi.txt | grep -q 200 || fail "Should have found /C/hi.txt"

echo "testing get"
ngxfs /hi.txt | grep -q helo || fail "Didn't find /hi.txt"

echo "testing ls"
ngxfs -ls / | grep -q "hi.txt" || fail "Incorrect ls list"

echo "testing del"
ngxfs -del /hi.txt
curl -I --silent $host/A/hi.txt | grep -q 404 || fail "Shouldn't have found /A/hi.txt"
curl -I --silent $host/B/hi.txt | grep -q 404 || fail "Shouldn't have found /B/hi.txt"
curl -I --silent $host/C/hi.txt | grep -q 404 || fail "Shouldn't have found /C/hi.txt"

echo "cleaning up"
curl -XDELETE http://0.0.0.0:8080/A/
curl -XDELETE http://0.0.0.0:8080/C/

echo "win"
