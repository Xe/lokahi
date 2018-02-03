#!/bin/sh

protoc --proto_path=$GOPATH/src:. \
	     --twirp_out=. \
	     --go_out=. \
	     ./lokahiadmin.proto
