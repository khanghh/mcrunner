package proto

//go:generate protoc --proto_path=../../proto --go-grpc_out=. --go_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative mcrunner.proto
