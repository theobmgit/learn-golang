all:
	protoc --go_out=proto/ proto/student/student.proto
	protoc --go-grpc_out=proto/ proto/student/student.proto

clean:
	rm gen/proto/*.go