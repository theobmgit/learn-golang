package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-student-management/docs"
	pb "go-student-management/proto/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	clientAddr = flag.String("c-addr", "localhost:50052", "The client address in the format of host:port")
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

// @title Student Management API
// @version 1.0
// @description Sample student management RPC-based API using gRPC server and gin client.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:50052
// @BasePath /api/v1
// @schemes http
func main() {
	const BasePath = "/api/v1"
	connection, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client error: %v\n", err)
	}
	defer connection.Close()

	client := pb.NewStudentServiceClient(connection)

	router := gin.Default()
	router.StaticFile("/favicon.ico", "resources/favicon.ico")

	v1 := router.Group(BasePath)
	{
		student := v1.Group("/students")
		{
			student.GET("", func(ctx *gin.Context) {
				getStudents(ctx, client)
			})
			student.GET("/:id", func(ctx *gin.Context) {
				getStudentById(ctx, client)
			})
			student.GET("/name/:name", func(ctx *gin.Context) {
				getStudentByName(ctx, client)
			})

			student.POST("", func(ctx *gin.Context) {
				addStudent(ctx, client)
			})

			student.PUT("/:id", func(ctx *gin.Context) {
				updateStudent(ctx, client)
			})

			student.DELETE("/:id", func(ctx *gin.Context) {
				deleteStudent(ctx, client)
			})
		}
	}

	docs.SwaggerInfo.BasePath = BasePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(*clientAddr); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

// getStudents godoc
// @Summary      Get all students
// @Description  Get all students' information
// @Tags         Student
// @Accept       json
// @Produce      json
// @Success      200  {array}  map[string][]string
// @Router       /students [get]
func getStudents(ctx *gin.Context, client pb.StudentServiceClient) {
	response, err := client.GetStudents(ctx, &pb.EmptyRequest{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// getStudentById godoc
// @Summary      Get student by student id
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID"
// @Success      200  {object}  map[string][]string
// @Router       /students/{id} [get]
func getStudentById(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")

	response, err := client.GetStudentById(ctx, &pb.StudentId{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// getStudentByName godoc
// @Summary      Get student by name
// @Description  Get student information based on student name
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        name path string true "Student name"
// @Success      200  {object}  map[string][]string
// @Router       /students/name/{name} [get]
func getStudentByName(ctx *gin.Context, client pb.StudentServiceClient) {
	name := ctx.Param("name")

	response, err := client.GetStudentByName(ctx, &pb.StudentName{Name: name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// addStudent godoc
// @Summary      Add new student
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        student body student.Student true "Student information"
// @Success      200  {object}  map[string][]string
// @Router       /students [post]
func addStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := client.AddStudent(ctx, student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// updateStudent godoc
// @Summary      Update student by ID
// @Description  Update student information based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to update"
// @Param        student body student.Student true "Student information to update"
// @Success      200  {object}  map[string][]string
// @Router       /students/{id} [put]
func updateStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	student.Id = ctx.Param("id")

	response, err := client.UpdateStudent(ctx, student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// deleteStudent godoc
// @Summary      Delete student by ID
// @Description  Delete student record based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to delete"
// @Success      200  {object}  map[string][]string
// @Router       /students/:id [delete]
func deleteStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")

	response, err := client.DeleteStudent(ctx, &pb.StudentId{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}
