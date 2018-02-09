#!/bin/sh

protoc --proto_path=$GOPATH/src:. \
	     --twirp_out=. \
	     --twirp_swagger_out=. \
	     --go_out=. \
	     ./lokahiadmin.proto
