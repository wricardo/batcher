#!/bin/bash
set -e

# Run test coverage on each subdirectories and merge the coverage profile.

echo "mode: count" > coverage.out
go test   -v -covermode=count -coverprofile=coverage.out

# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test  -v -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> coverage.out
        rm $dir/profile.tmp
    fi
fi
done
