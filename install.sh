#!/bin/bash

#let script exit when command fails ;)
set -e
#let script exit when script uses undeclared variables
set -u

user=$(whoami)
echo "Welcome to the installation," $user!
echo ""

if ! [ -x "$(command -v go)" ]; then
  echo "["`date`"][ERROR]: go is not installed." >&2
  exit 1
else
  echo "["`date`"][SUCCESS]: Go installed."
fi

PACK=$(echo $GOPATH)
if [[ -d $PACK/src/github.com/gocolly/colly ]]; then
    echo "["`date`"][SUCCESS]: Colly package is found, great."
else
    echo "["`date`"][WARNING]: Colly package is not found, installing.."
    go get github.com/gocolly/colly
    echo "["`date`"][INFO]: Colly package done..1 more"
    go get github.com/gocolly/colly/extensions
    echo "["`date`"][INFO]: Colly extension package done.."
fi

echo ""
echo "["`date`"][INFO]: Building the app.."
go run .
echo "["`date`"][INFO]: OUTPUT: [zakon.csv]"
