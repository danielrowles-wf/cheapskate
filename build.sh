#! /bin/bash

set -e

PREFIX=github.com/danielrowles-wf/cheapskate/gen-go/
IMPORT="github.com/Workiva/frugal/lib/go"
TOPDIR=$PWD

echo "Building parsimony"
go get github.com/Workiva/parsimony

echo "Building frugal"
go get github.com/Workiva/frugal

if [ -e ./gen-go ]; then
    echo "Remove existing gen-go directory"
    rm -Rf gen-go
fi

if [ -e ./stage ]; then
    echo "Remove existing stage directory"
    rm -Rf stage
fi

echo "Fetch all required IDL files"
$GOPATH/bin/parsimony --staging stage stingy.frugal

echo "Generate GO code"
$GOPATH/bin/frugal --gen=go:package_prefix=$PREFIX -r stage/stingy.frugal

for dir in cheapskate client; do
    echo "Building <$dir>"
    cd $TOPDIR/$dir
    if [ -e $dir ]; then
        rm $dir
    fi
    go build .
    cd $TOPDIR
done
