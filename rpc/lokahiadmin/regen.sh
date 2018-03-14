#!/bin/sh

protoc --proto_path=$GOPATH/src:. \
	     --twirp_out=. \
	     --twirp_swagger_out=. \
	     --nrpc_out=. \
	     --go_out=. \
	     ./lokahiadmin.proto

# XXX workaround to prevent
# # github.com/Xe/lokahi/rpc/lokahiadmin
# rpc/lokahiadmin/lokahiadmin.nrpc.go:98: Println call has possible formatting directive %s
# rpc/lokahiadmin/lokahiadmin.nrpc.go:217: Println call has possible formatting directive %s
# rpc/lokahiadmin/lokahiadmin.nrpc.go:336: Println call has possible formatting directive %s

sed -i 's/Println/Printf/g' lokahiadmin.nrpc.go
