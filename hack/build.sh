#!/bin/bash

# This script builds the project, without the need to understand how to setup a GOPATH
# inspired from docker's script: https://github.com/docker/docker/blob/master/hack/make.sh

export CODE_BASE='../lambdaManager'

export LAMBDA_PACKAGE='src/github.com/tylerharter/open-lambda'
export GOPATH='.gopath'

export WORKER='lambdaManager/server'
export CLIENT='lambdaManager/prof/client'

mkdir -p ${GOPATH}/${LAMBDA_PACKAGE}
ln -sf $(realpath ${CODE_BASE}) ${GOPATH}/${LAMBDA_PACKAGE}

# now that gopath exists, we re-export as absolute, required by go
export GOPATH=$(realpath ${GOPATH})

cd ${GOPATH}/${LAMBDA_PACKAGE}/${WORKER} && go get && go build
cd ${GOPATH}/${LAMBDA_PACKAGE}/${CLIENT} && go get && go build
