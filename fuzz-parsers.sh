#!/usr/bin/env sh

GOPKG=./syntax/languages
FUZZTIME=15m

go test $GOPKG -list Fuzz | grep Fuzz | while read line; do
    echo "======== $line =========="
    go test $GOPKG -fuzztime $FUZZTIME -fuzz "$line"
done
