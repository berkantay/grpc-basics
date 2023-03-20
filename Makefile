create:
	protoc --proto_path=grpc/proto/ ./grpc/proto/*.proto --go_out=./grpc/proto/.
	protoc --proto_path=grpc/proto/ ./grpc/proto/*.proto --go-grpc_out=./grpc/proto/.

clean:
	rm grpc/proto/*.go