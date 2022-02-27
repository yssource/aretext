#!/usr/bin/env sh

GOPKG=./syntax/languages
FUZZTIME=5s

go test $GOPKG -list Fuzz | grep Fuzz | while read line; do
    echo "======== $line =========="
    go test $GOPKG -fuzztime $FUZZTIME -fuzz "$line"
done
