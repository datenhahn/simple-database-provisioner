#!/bin/bash

set -e

APPVERSION=v0.1.3

docker build . -t ecodia/simple-database-provisioner

docker tag ecodia/simple-database-provisioner ecodia/simple-database-provisioner:$APPVERSION
