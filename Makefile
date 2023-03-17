create:
	protoc --proto_path=internal/adapters/driving/proto internal/adapters/driving/proto/*.proto --go_out=internal/adapters/driving/proto/.
	protoc --proto_path=internal/adapters/driving/proto internal/adapters/driving/proto/*.proto --go-grpc_out=internal/adapters/driving/proto/.

clean:
	rm internal/adapters/driving/proto/*.go