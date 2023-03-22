build:
	go build -o user-management-service -ldflags="-X 'main.Version=v1.0.0'" cmd/main.go
create:
	protoc --proto_path=grpc/proto/ ./grpc/proto/*.proto --go_out=./grpc/proto/.
	protoc --proto_path=grpc/proto/ ./grpc/proto/*.proto --go-grpc_out=./grpc/proto/.
clean:
	rm grpc/proto/*.go
test:
	go test -v ./... --cover