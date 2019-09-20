#!/usr/bin/env bash

set -e
cd $GOPATH/src/aether-core
echo "Running all tests and generating coverage profile for the entire project. It will be shown in browser once complete."
for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
echo "mode: atomic" > result.txt && cat coverage.txt | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> result.txt && go tool cover -html=result.txt && rm result.txt && rm coverage.txt