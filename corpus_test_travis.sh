#!/bin/bash
set -ev

if [ "${TRAVIS_PULL_REQUEST}" = "false" ] && [[ `go version` = *"go1.12"* ]]; then
	go test -v -race -run TestCorpus -randtests 50 -corpus .travis-corpus -gencorpus .
	cd .travis-corpus && sha256sum -c --quiet --strict ../corpus-sha256sums
fi