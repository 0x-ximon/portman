package main

// Create the proto directory if it doesn't exist
//go:generate mkdir -p proto

// Generate the code relative to the output directory
//go:generate protoc -I=../.proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative ../.proto/orders.proto
