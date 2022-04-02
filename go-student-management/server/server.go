package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	pb "go-student-management/proto/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	jsonDBFile = flag.String("json_db_file", "database/db.json", "A json file containing a list of features")
	port       = flag.Int("port", 50051, "The server port")
)

type studentServiceServer struct {
	students []*pb.Student
	pb.UnimplementedStudentServiceServer
}

func newStudentServiceServer() *studentServiceServer {
	s := &studentServiceServer{}
	s.loadStudents(*jsonDBFile)
	return s
}

func (s *studentServiceServer) loadStudents(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	data, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(data, &s.students); err != nil {
		log.Fatalf("Failed to load default students: %v", err)
	}
}

func (s *studentServiceServer) searchById(studentId string) (bool, int) {
	for index, student := range s.students {
		if student.Id == studentId {
			return true, index
		}
	}
	return false, -1
}

func (s *studentServiceServer) remove(index int) {
	a := s.students
	a[index] = a[len(a)-1]
	a[len(a)-1] = &pb.Student{}
	a = a[:len(a)-1]
}

func (s *studentServiceServer) GetStudentById(_ context.Context, studentId *pb.StudentId) (*pb.Student, error) {
	if ok, index := s.searchById(studentId.Id); ok {
		student := s.students[index]
		return &pb.Student{Id: student.Id, Name: student.Name, Age: student.Age, Class: student.Class, Cpa: student.Cpa}, nil
	}
	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", studentId.Id)
}

func (s *studentServiceServer) GetStudentByName(_ context.Context, studentName *pb.StudentName) (*pb.Student, error) {
	for _, student := range s.students {
		if strings.Contains(student.Name, studentName.Name) {
			return &pb.Student{Id: student.Id, Name: student.Name, Age: student.Age, Class: student.Class, Cpa: student.Cpa}, nil
		}
	}

	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with name=%s", studentName.Name)
}

func (s *studentServiceServer) GetStudents(context.Context, *pb.EmptyRequest) (*pb.Students, error) {
	return &pb.Students{Students: s.students}, nil
}

func (s *studentServiceServer) AddStudent(_ context.Context, student *pb.Student) (*pb.Student, error) {
	if ok, _ := s.searchById(student.Id); ok {
		return &pb.Student{}, status.Errorf(codes.AlreadyExists, "Student with id %v already exist!", student.Id)
	}

	s.students = append(s.students, student)
	return &pb.Student{Id: student.Id, Name: student.Name, Age: student.Age, Class: student.Class, Cpa: student.Cpa}, nil
}

func (s *studentServiceServer) UpdateStudent(_ context.Context, newStudent *pb.Student) (*pb.Student, error) {
	if ok, index := s.searchById(newStudent.Id); ok {
		s.students[index] = newStudent
		return &pb.Student{Id: newStudent.Id, Name: newStudent.Name, Age: newStudent.Age, Class: newStudent.Class, Cpa: newStudent.Cpa}, nil
	}

	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", newStudent.Id)
}

func (s *studentServiceServer) DeleteStudent(_ context.Context, studentId *pb.StudentId) (*pb.Student, error) {
	if ok, index := s.searchById(studentId.Id); ok {
		s.remove(index)
		student := s.students[index]
		return &pb.Student{Id: student.Id, Name: student.Name, Age: student.Age, Class: student.Class, Cpa: student.Cpa}, nil
	}

	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", studentId.Id)
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
	}
	log.Printf("Server successfully load on %v\n", listener.Addr())

	grpcServer := grpc.NewServer()
	pb.RegisterStudentServiceServer(grpcServer, newStudentServiceServer())
	err2 := grpcServer.Serve(listener)
	if err2 != nil {
		log.Fatalf("Failed to serve %v\n", err2)
	}
}
