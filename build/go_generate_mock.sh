#!/bin/bash

WITH_DOCKER=

run_mockery() {
	if [ -z "$WITH_DOCKER" ]; then 
		mockery "$@"
	else
		docker run -it --rm -v "$PWD:/src" -w /src vektra/mockery:v2.8 "$@"
	fi
}

#run_mockery --name=Writer --srcpkg=io --output=test/mock/mockio --outpkg=mockio --filename=mockio.go
run_mockery --name=Logger --structname MockLogger --dir=. --output=. --outpkg=justlog --filename=logger_mock.go
