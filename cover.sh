#!/usr/bin/env bash

set -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

usage()
{
cat << EOF
Usage: $(basename $0) --mode CoverMode

This script will run 'go test' on all non-vendor subdirs, then 'go tool cover'
for test coverage. Any 'go test' failures will exit the script. If all tests
are successful, it outputs a 'coverage.html' file with test coverage result.

OPTIONS:
 -h                 Show this message
 -m or --mode       Test coverage mode: set, count, or atomic (default)
EOF
}

# Set default variables
coverMode="atomic"
while :
do
    case $1 in
        -h | --help | -\?)
        usage
        exit 0      # This not an error, User asked help. Don't do "exit 1"
        ;;
        -m | --mode)
            coverMode=$2
            shift 2
            ;;
        *)
            break
            ;;
    esac
done

echo "mode: $coverMode" | tee coverage.txt

cd $DIR
for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=profile.out -covermode=$coverMode $d
    if [ -f profile.out ]; then
        tail +2 profile.out >> coverage.txt
        rm profile.out
    fi
done
go tool cover -html=coverage.txt -o cover.html
