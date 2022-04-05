package main

import (
	"flag"
	"fmt"
	pb "go-student-management/proto/student"
	studentServer "go-student-management/server/student"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
	}
	log.Printf("Server successfully load on %v\n", listener.Addr())

	grpcServer := grpc.NewServer()
	pb.RegisterStudentServiceServer(grpcServer, studentServer.NewStudentServiceServer())
	err2 := grpcServer.Serve(listener)
	if err2 != nil {
		log.Fatalf("Failed to serve %v\n", err2)
	}
}
