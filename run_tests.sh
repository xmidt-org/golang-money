#!/usr/bin/env bash
set -ex

# Script does performs two task:
# (1) Automatic checks using: gofmt, go vet, gosimple, and uncovert 
# (2) Runs go tests 

# gometalinter (github.com/alecthomas/gometalinter) is used to run each each
# static checker.

testgmoney () { 

    # (1) 
    gometalinter --vendor --disable-all --deadline=10m \
        --enable=gofmt \
        --enable=vet \
        --enable=gosimple \
        --enable=unconvert \
        --enable=ineffassign \
        ./...
    if [ $? != 0 ]; then 
        echo 'Gometalinter has some complaints'
        exit 1
    fi 

    # (2) 
    env GORACE= 'halt_on_errors=1' go test -coverprofile=coverage.txt ./...
    echo "running go tests"

    if [ $? != 0 ]; then 
        echo 'GO tests failed'
        exit 1
    fi

    echo "----------------------"
    echo "Tests completed successfully!"
}

testgmoney


