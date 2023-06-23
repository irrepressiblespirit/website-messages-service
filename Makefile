
default: help

mongohost?=172.16.0.16:27017
action?=up

## make clean - clean app
client:
	go build -o bin/client cmd/client/*.go

## make run - run application
run: 
	go run cmd/main.go

## make runlocal - run application with local config
runlocal:
	config=config.local.yaml go run cmd/main.go

## make proto - generate Proto files
proto:
	protoc -I ./ api/proto/api.proto --go_out=pkg/grpcapi --go-grpc_out=pkg/grpcapi
	protoc -I ./ api/centrifugo/api.proto --go_out=pkg/centrifugopb --go-grpc_out=pkg/centrifugopb

## make evans - connect to local gRPC server
## for user auth use: header token=develop
evans:
	evans api/proto/api.proto -p 8081

## make grpcui-local - run UI for local gRPC server
grpcui-local:
	grpcui -proto api/proto/api.proto -rpc-header 'token: develop' -plaintext 127.0.0.1:8081