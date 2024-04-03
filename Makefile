.PHONY: all clean generate run

all: generate run

generate: clean
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. protos/location.proto

run:
	go run .
