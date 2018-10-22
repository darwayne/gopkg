#!/usr/bin/env bash

cd ../
./scripts/build.sh -p
cd docker
./docker.sh down
./docker.sh up -d --build

