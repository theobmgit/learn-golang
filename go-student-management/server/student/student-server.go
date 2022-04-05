package student

import (
	"context"
	"encoding/json"
	"flag"
	"go-student-management/proto/student"
	pb "go-student-management/proto/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	jsonDBFile = flag.String("json_db_file", "database/db.json", "A json file containing a list of features")
)

type studentServiceServer struct {
	students []*student.Student
	student.UnimplementedStudentServiceServer
}

func NewStudentServiceServer() *studentServiceServer {
	s := &studentServiceServer{}
	s.loadStudents(*jsonDBFile)
	return s
}

func (s *studentServiceServer) loadStudents(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(data, &s.students); err != nil {
		log.Fatalf("Failed to load default students: %v", err)
	}
}

func (s *studentServiceServer) saveStudents(filePath string) {
	dataJson, err := json.MarshalIndent(s.students, "", "\t")
	if err != nil {
		log.Fatalf("Failed to load default students: %v", err)
	}

	if err2 := os.WriteFile(filePath, dataJson, 0644); err2 != nil {
		log.Fatalf("Failed to save students: %v", err2)
	}
}

func (s *studentServiceServer) searchById(studentId string) (bool, int) {
	for index, stud := range s.students {
		if stud.Id == studentId {
			return true, index
		}
	}
	return false, -1
}

func (s *studentServiceServer) remove(index int) {
	a := s.students
	a[index] = a[len(a)-1]
	a[len(a)-1] = &student.Student{}
	a = a[:len(a)-1]
}

func (s *studentServiceServer) GetStudentById(_ context.Context, studentId *student.StudentId) (*student.Student, error) {
	if ok, index := s.searchById(studentId.Id); ok {
		stud := s.students[index]
		return &pb.Student{Id: stud.Id, Name: stud.Name, Age: stud.Age, Class: stud.Class, Cpa: stud.Cpa}, nil
	}
	return &student.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", studentId.Id)
}

func (s *studentServiceServer) GetStudentsByName(_ context.Context, studentName *student.StudentName) (*student.Students, error) {
	results := pb.Students{}
	for _, stud := range s.students {
		if strings.Contains(stud.Name, studentName.Name) {
			results.Students = append(results.Students, &pb.Student{Id: stud.Id, Name: stud.Name, Age: stud.Age, Class: stud.Class, Cpa: stud.Cpa})
		}
	}

	if len(results.Students) != 0 {
		return &results, nil
	} else {
		return &results, status.Errorf(codes.NotFound, "Do not found stud with name=%s", studentName.Name)
	}
}

func (s *studentServiceServer) GetStudents(context.Context, *student.EmptyRequest) (*student.Students, error) {
	return &student.Students{Students: s.students}, nil
}

func (s *studentServiceServer) AddStudent(_ context.Context, stud *student.Student) (*student.Student, error) {
	if ok, _ := s.searchById(stud.Id); ok {
		return &pb.Student{}, status.Errorf(codes.AlreadyExists, "Student with id %v already exist!", stud.Id)
	}

	s.students = append(s.students, stud)
	s.saveStudents(*jsonDBFile)
	return &pb.Student{Id: stud.Id, Name: stud.Name, Age: stud.Age, Class: stud.Class, Cpa: stud.Cpa}, nil
}

func (s *studentServiceServer) UpdateStudent(_ context.Context, newStudent *student.Student) (*student.Student, error) {
	if ok, index := s.searchById(newStudent.Id); ok {
		s.students[index] = newStudent
		s.saveStudents(*jsonDBFile)
		return &pb.Student{Id: newStudent.Id, Name: newStudent.Name, Age: newStudent.Age, Class: newStudent.Class, Cpa: newStudent.Cpa}, nil
	}

	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", newStudent.Id)
}

func (s *studentServiceServer) DeleteStudent(_ context.Context, studentId *student.StudentId) (*student.Student, error) {
	if ok, index := s.searchById(studentId.Id); ok {
		s.remove(index)
		stud := s.students[index]
		s.saveStudents(*jsonDBFile)
		return &pb.Student{Id: stud.Id, Name: stud.Name, Age: stud.Age, Class: stud.Class, Cpa: stud.Cpa}, nil
	}

	return &pb.Student{}, status.Errorf(codes.NotFound, "Do not found student with id=%s", studentId.Id)
}
