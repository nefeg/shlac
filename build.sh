#!/usr/bin/env bash

GOPATH=$(pwd)

export GOPATH=$GOPATH
export GOBIN=$GOPATH/bin

go install shlac shlacd; # SHlack Like As Cron

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            echo "HELP!"
            exit 0
            ;;

        -r|--run)
            `$GOBIN/shlacd $GOPATH/config.json`
            ;;
    esac
done
#echo "Done"