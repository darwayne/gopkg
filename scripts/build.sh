#!/usr/bin/env bash

usage() {
  echo "Usage: $0 [-vhnp]" 1>&2
  echo -e "\t-v:\tVerbose"
  echo -e "\t-h:\tShow usage"
  echo -e "\t-n:\tNo build. Useful if your IDE auto builds and runs. Only installs dependencies"
  echo -e "\t-p:\tBuild For Prod"
  exit 1
}

## Download dep for dep resolution
if [ -z "$(which dep)" ]; then
  echo Installing dep...
  go get -u github.com/golang/dep/cmd/dep
  echo Done.
fi

## Download dep for dep resolution
if [ -z "$(which packr)" ]; then
  echo Installing packr...
  go get -u github.com/gobuffalo/packr/...
  echo Done.
fi

while getopts ":vhnp" o; do
    case "${o}" in
        v)
            GOFLAGS="-x"
            ;;
        h)
            usage
            ;;
        n)
            NO_COMPILE="1"
            ;;
        p)  PRODUCTION="1"
            ;;
        *)
            echo Unknown flag: ${OPTARG}
            usage
            ;;
    esac
done
shift $((OPTIND-1))

## update deps
VENDOR=vendor/
if [ -d $VENDOR ]; then
  rm -rf $VENDOR
fi

dep ensure

# run packr binary so files will be generated
packr

if [ -z "$NO_COMPILE" ]; then

   if [ -z "$PRODUCTION" ]; then
        go build ${GOFLAGS} -a \
            -installsuffix cgo \
            github.com/darwayne/gopkg/cmd/gopkg
   else
      CGO_ENABLED=0 GOOS=linux go build ${GOFLAGS} -a \
          -installsuffix cgo \
          github.com/darwayne/gopkg/cmd/gopkg
   fi

fi
