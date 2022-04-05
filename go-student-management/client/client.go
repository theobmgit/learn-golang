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
	"math"
	"net/http"
	"regexp"
	"strings"
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
				GetStudents(ctx, client)
			})
			student.GET("/:id", func(ctx *gin.Context) {
				GetStudentById(ctx, client)
			})
			student.GET("/name/:name", func(ctx *gin.Context) {
				GetStudentByName(ctx, client)
			})

			student.POST("", func(ctx *gin.Context) {
				AddStudent(ctx, client)
			})

			student.PUT("/:id", func(ctx *gin.Context) {
				UpdateStudent(ctx, client)
			})

			student.DELETE("/:id", func(ctx *gin.Context) {
				DeleteStudent(ctx, client)
			})
		}
	}

	docs.SwaggerInfo.BasePath = BasePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(*clientAddr); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

// GetStudents godoc
// @Summary      Get all students
// @Description  Get all students' information
// @Tags         Student
// @Accept       json
// @Produce      json
// @Success      200  {array}  map[string][]string
// @Router       /students [get]
func GetStudents(ctx *gin.Context, client pb.StudentServiceClient) {
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

// GetStudentById godoc
// @Summary      Get student by student id
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID" minLength(8)
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/{id} [get]
func GetStudentById(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")
	if validateId(id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	}

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

// GetStudentByName godoc
// @Summary      Get student by name
// @Description  Get student information based on student name
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        name path string true "Student name" pattern("^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$)
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/name/{name} [get]
func GetStudentByName(ctx *gin.Context, client pb.StudentServiceClient) {
	name := ctx.Param("name")
	if validateName(name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	}
	preprocessName(&name)

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

// UpdateStudent godoc
// @Summary      Update student by ID
// @Description  Update student information based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to update" minLength(8)
// @Param        student body student.Student true "Student information to update"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/{id} [put]
func UpdateStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	student.Id = ctx.Param("id")

	if validateId(student.Id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	} else if validateName(student.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	} else if validateAge(student.Age) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Age must be an integer larger than 0",
		})
		return
	} else if validateClass(student.Class) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class",
		})
		return
	} else if validateCpa(student.Cpa) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "CPA must be between 0.0 and 4.0",
		})
		return
	}
	preprocessName(&student.Name)
	preprocessCpa(&student.Cpa)
	preprocessClass(&student.Class)

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

// DeleteStudent godoc
// @Summary      Delete student by ID
// @Description  Delete student record based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to delete" minLength(8)
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/:id [delete]
func DeleteStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")

	if validateId(id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	}

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

// AddStudent godoc
// @Summary      Add new student
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        student body student.Student true "Student information"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students [post]
func AddStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if validateId(student.Id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	} else if validateName(student.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	} else if validateAge(student.Age) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Age must be an integer larger than 0",
		})
		return
	} else if validateClass(student.Class) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class",
		})
		return
	} else if validateCpa(student.Cpa) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "CPA must be between 0.0 and 4.0",
		})
		return
	}
	preprocessName(&student.Name)
	preprocessCpa(&student.Cpa)
	preprocessClass(&student.Class)

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

func validateId(id string) bool {
	return len(id) != 8
}

func validateName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$", name)
	return match
}

func validateAge(age int32) bool {
	return age >= 0
}

func validateClass(class string) bool {
	return true
}

func validateCpa(Cpa float32) bool {
	return Cpa >= 0.0 && Cpa <= 4.0
}

func preprocessName(name *string) {
	*name = strings.ToTitle(strings.ToLower(*name))
}

func preprocessClass(class *string) {
	*class = strings.ToUpper(*class)
}

func preprocessCpa(Cpa *float32) {
	*Cpa = float32(math.Ceil(float64(*Cpa*100)) / 100)
}
