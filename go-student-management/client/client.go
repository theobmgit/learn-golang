package main

import (
	"context"
	"flag"
	pb "go-student-management/proto/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	connection, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client error: %v\n", err)
	}
	defer connection.Close()

	client := pb.NewStudentServiceClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request, err := client.GetStudentById(ctx, &pb.StudentId{Id: "20124721"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	request, err = client.GetStudentByName(ctx, &pb.StudentName{Name: "Bui"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	request2, err := client.GetStudents(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request2.GetStudents())

	request, err = client.AddStudent(ctx, &pb.Student{Id: "20192642", Name: "Nguyen Luu Hoang Minh", Age: 20, Class: "ICT-01", Cpa: 3.6})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	request2, err = client.GetStudents(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request2.GetStudents())

	request, err = client.UpdateStudent(ctx, &pb.Student{Id: "20192642", Name: "Nguyen Luu Hoang Minh", Age: 20, Class: "ICT-01", Cpa: 4.0})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	request, err = client.GetStudentById(ctx, &pb.StudentId{Id: "20192642"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	log.Println("Delete")
	request, err = client.DeleteStudent(ctx, &pb.StudentId{Id: "20192642"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request)

	request2, err = client.GetStudents(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student: %v\n", request2.GetStudents())
}
