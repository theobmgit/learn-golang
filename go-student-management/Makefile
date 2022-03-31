all:
	protoc --go_out=proto/ proto/student.proto
	protoc --go-grpc_out=proto/ proto/student.proto

clean:
	rm gen/proto/*.go