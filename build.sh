#!/bin/sh

SVN_VERSION=`git rev-parse --short HEAD || echo "GitNotFound"`
APP_VERSION=`head -n 3 update.sh|grep LAST_VERSION|awk -F"=" '{print $2}'`

go build -ldflags  "-X main.buildTime=`date  +%Y%m%d-%H%M%S` -X main.binaryVersion=$APP_VERSION -X main.svnRevision=${SVN_VERSION}"

