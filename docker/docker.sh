#!/usr/bin/env bash

# Name space current folder appropriately in docker
# This avoids conflicts with other docker projects in
# a similar folder structure
docker-compose --project-name "gopkg" $@
