#!/usr/bin/env bash

docker exec -it $(docker ps | grep "gopkg_app" |  cut -d' ' -f1) /bin/gopkg -db.host db -db.store true